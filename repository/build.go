package repository

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"encoding/json"
	"bytes"

	"github.com/opencontainers/image-spec/identity"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/reference"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
	"github.com/containerd/containerd/diff"
	"github.com/pauldotknopf/darch/workspace"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/content"
	digest "github.com/opencontainers/go-digest"
)

// BuildRecipe Builds a recipe.
func (session *Session) BuildRecipe(ctx context.Context, recipe recipes.Recipe, tag string, buildPrefix string, environmentVariables map[string]string) error {

	ctx = namespaces.WithNamespace(ctx, "darch")

	if len(tag) == 0 {
		tag = "local"
	}

	inheritsRef, err := reference.Parse(recipe.Inherits)
	if err != nil {
		return err
	}

	// If inherited image defines no tag, use the tag we are building with
	if len(inheritsRef.Object) == 0 {
		inheritsRef.Object = tag
	}

	img, err := session.client.GetImage(ctx, inheritsRef.String())
	if err != nil {
		// maybe it was because we don't have it? let's try to fetch it
		if img, err = session.Pull(ctx, inheritsRef.String()); err != nil {
			return err
		}
	}

	ws, err := workspace.NewWorkspace("/tmp")
	if err != nil {
		return err
	}
	defer workspace.DestroyWorkspace(ws)

	mounts := []specs.Mount{
		specs.Mount{
			Destination: "/recipes",
			Type:        "bind",
			Source:      recipe.RecipesDir,
			Options:     []string{"rbind", "ro"},
		},
	}

	if utils.FileExists("/etc/resolv.conf") {
		utils.CopyFile("/etc/resolv.conf", path.Join(ws.Path, "resolv.conf"))
		mounts = append(mounts, specs.Mount{
			Destination: "/etc/resolv.conf",
			Type:        "bind",
			Source:      path.Join(ws.Path, "resolv.conf"),
			Options:     []string{"rbind", "rw"},
		})
	}

	// Let's create the snapshot that all of our containers will run off of
	snapshotKey := utils.NewID()
	session.createSnapshot(ctx, snapshotKey, img)
	defer session.deleteSnapshot(ctx, snapshotKey)

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", "/darch-prepare"),
			),
		},
	}); err != nil {
		return err
	}

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", fmt.Sprintf("/darch-runrecipe %s", recipe.Name)),
			),
		},
	}); err != nil {
		return err
	}

	if err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", "/darch-teardown"),
			),
		},
	}); err != nil {
		return err
	}

	return session.createImageFromSnapshot(ctx, img, snapshotKey, recipe.Name + ":" + tag)
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

func (session *Session) createImageFromSnapshot(ctx context.Context, img containerd.Image, activeSnapshotKey string, newReference string) error {

	contentStore := session.client.ContentStore()
	snapshotService := session.client.SnapshotService(containerd.DefaultSnapshotter)
	imgTarget := img.Target()

	// First, let's get the parent image digest, so that we can
	// later create a new one from it, with a new layer added to it.
	p, err := content.ReadBlob(ctx, contentStore, imgTarget.Digest)
	if err != nil {
		return err
	}
	var manifest ocispec.Manifest
	if err := json.Unmarshal(p, &manifest); err != nil {
		return err
	}

	snapshot, err := snapshotService.Stat(ctx, activeSnapshotKey)
	if err != nil {
		return err
	}

	upperMounts, err := snapshotService.Mounts(ctx, activeSnapshotKey)
	if err != nil {
		return err
	}

	lowerMounts, err := snapshotService.View(ctx, "temp-readonly-parent", snapshot.Parent)
	if err != nil {
		return err
	}
	defer snapshotService.Remove(ctx, "temp-readonly-parent")

	diffs, err := session.client.DiffService().DiffMounts(ctx,
		lowerMounts,
		upperMounts,
		diff.WithMediaType(ocispec.MediaTypeImageLayerGzip),
		diff.WithReference("custom-ref"))
	if err != nil {
		return err
	}
	// Add our new layer to the image manifest
	for _,layer := range manifest.Layers {
		fmt.Println(layer.Digest)
		fmt.Println(layer.MediaType)
	}
	manifest.Layers = append(manifest.Layers, diffs)
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	manifestDigest := digest.FromBytes(manifestBytes)
	if err := content.WriteBlob(ctx,
		contentStore,
		"ref1",
		bytes.NewReader(manifestBytes),
		int64(len(manifestBytes)),
		manifestDigest); err != nil {
			return err
	}

	_, err = session.client.ImageService().Create(ctx,
		images.Image{
			Name:   newReference,
			Target: ocispec.Descriptor{
				Digest:    manifestDigest,
				Size:      int64(len(manifestBytes)),
				MediaType: ocispec.MediaTypeImageManifest,
			},
		})
	if err != nil {
		return err
	}

	return nil
}