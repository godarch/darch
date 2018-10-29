package repository

import (
	"context"
	"fmt"
	"runtime"

	"github.com/opencontainers/image-spec/identity"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/godarch/darch/pkg/recipes"
	"github.com/godarch/darch/pkg/reference"
	"github.com/godarch/darch/pkg/repository/manifest"
	"github.com/godarch/darch/pkg/utils"
	"github.com/godarch/darch/pkg/workspace"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// BuildRecipe Builds a recipe.
func (session *Session) BuildRecipe(ctx context.Context, recipe recipes.Recipe, tag string, imagePrefix string, env []string) (reference.ImageRef, error) {

	ctx = namespaces.WithNamespace(ctx, "darch")

	if len(tag) == 0 {
		tag = "latest"
	}

	newImage, err := reference.ParseImage(imagePrefix + recipe.Name + ":" + tag)
	if err != nil {
		return nil, err
	}

	// Use the image prefix when inheriting local recipes.
	// External references are expected to be fully qualified.
	inherits := recipe.Inherits
	if !recipe.InheritsExternal {
		inherits = imagePrefix + inherits
	}

	// NOTE: We use ParseImageWithDefaultTag here.
	// This allows recipes to use specific tags, but when
	// they aren't, it uses the tag the we are building
	// the recipe with.
	// This allows use to "darch build -t custom-tag base base-common"
	// and each built image will use the appropriate inherited image.
	inheritsRef, err := reference.ParseImageWithDefaultTag(inherits, newImage.Tag())
	if err != nil {
		return newImage, err
	}

	img, err := session.client.GetImage(ctx, inheritsRef.FullName())
	if err != nil {
		return newImage, err
	}

	ws, err := workspace.NewWorkspace("/tmp")
	if err != nil {
		return newImage, err
	}
	defer ws.Destroy()

	mounts, err := createTempMounts(ws.Path)

	mounts = append(mounts, specs.Mount{
		Destination: "/recipes",
		Type:        "bind",
		Source:      recipe.RecipesDir,
		Options:     []string{"rbind", "ro"},
	})

	// Prevent garbage collection while we work.
	ctx, done, err := session.client.WithLease(ctx)
	if err != nil {
		return newImage, err
	}
	defer done(ctx)

	// Let's create the snapshot that all of our containers will run off of
	snapshotKey := utils.NewID()
	err = session.createSnapshot(ctx, snapshotKey, img)
	if err != nil {
		return newImage, err
	}
	defer session.deleteSnapshot(ctx, snapshotKey)

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithEnv(env),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", "/darch-prepare"),
			),
		},
	}); err != nil {
		return newImage, err
	}

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithEnv(env),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", fmt.Sprintf("/darch-runrecipe %s", recipe.Name)),
			),
		},
	}); err != nil {
		return newImage, err
	}

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithEnv(env),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", "/darch-teardown"),
			),
		},
	}); err != nil {
		return newImage, err
	}

	return newImage, session.createImageFromSnapshot(ctx, img, snapshotKey, newImage)
}

func (session *Session) createSnapshot(ctx context.Context, snapshotKey string, img containerd.Image) error {
	diffIDs, err := img.RootFS(ctx)
	if err != nil {
		return err
	}
	parent := identity.ChainID(diffIDs).String()
	if _, err := session.client.SnapshotService(containerd.DefaultSnapshotter).Prepare(ctx, snapshotKey, parent); err != nil {
		return err
	}
	return nil
}

func (session *Session) deleteSnapshot(ctx context.Context, snapshotKey string) error {
	return session.client.SnapshotService(containerd.DefaultSnapshotter).Remove(ctx, snapshotKey)
}

func (session *Session) createImageFromSnapshot(ctx context.Context, img containerd.Image, activeSnapshotKey string, newImage reference.ImageRef) error {
	// First, let's get the parent image manifest so that we can
	// later create a new one from it, with a new layer added to it.
	m, err := manifest.LoadManifest(ctx, session.content, img.Target())
	if err != nil {
		return err
	}

	snapshot, err := session.snapshotter.Stat(ctx, activeSnapshotKey)
	if err != nil {
		return err
	}

	upperMounts, err := session.snapshotter.Mounts(ctx, activeSnapshotKey)
	if err != nil {
		return err
	}

	lowerMounts, err := session.snapshotter.View(ctx, "temp-readonly-parent", snapshot.Parent)
	if err != nil {
		return err
	}
	defer session.snapshotter.Remove(ctx, "temp-readonly-parent")

	// Generate a diff in content store
	diffs, err := session.client.DiffService().Compare(ctx,
		lowerMounts,
		upperMounts,
		diff.WithMediaType(ocispec.MediaTypeImageLayerGzip),
		diff.WithReference("custom-ref"))
	if err != nil {
		return err
	}

	// Add our new layer to the image manifest
	err = m.AddLayer(ctx, session.content, diffs)

	// Let's see if the image exists already, if so, let's delete it
	_, err = session.client.GetImage(ctx, newImage.FullName())
	if err == nil {
		session.client.ImageService().Delete(ctx, newImage.FullName(), images.SynchronousDelete())
	}

	_, err = session.client.ImageService().Create(ctx,
		images.Image{
			Name: newImage.FullName(),
			Target: ocispec.Descriptor{
				Digest:    m.Descriptor().Digest,
				Size:      m.Descriptor().Size,
				MediaType: m.Descriptor().MediaType,
			},
		})
	if err != nil {
		return err
	}

	// This will create the required snapshot for the new layer,
	// which will allow us to run the image immediately.
	imageBuilt, err := session.client.GetImage(ctx, newImage.FullName())
	if err != nil {
		return err
	}
	err = imageBuilt.Unpack(ctx, containerd.DefaultSnapshotter)
	if err != nil {
		return err
	}

	return nil
}
