package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"runtime"

	"github.com/opencontainers/image-spec/identity"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/reference"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
	"github.com/pauldotknopf/darch/workspace"
)

const containerdUncompressed = "containerd.io/uncompressed"

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

	return session.createImageFromSnapshot(ctx, img, snapshotKey, recipe.Name+":"+tag)
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

func (session *Session) patchImageConfig(ctx context.Context, ref string, manifest *ocispec.Manifest, newLayerDigest digest.Digest) error {
	// Get the current image configuration.
	p, err := content.ReadBlob(ctx, session.client.ContentStore(), manifest.Config.Digest)
	if err != nil {
		return err
	}

	// Deserialize the image configuration to a generic json object.
	// We do this so that we can patch it, without requiring knowledge
	// of the entire schema.
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(p, &m); err != nil {
		return err
	}

	// Pull the rootfs section out, so that we can append a layer to the diff_ids array.
	var rootFS ocispec.RootFS
	p, err = m["rootfs"].MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(p, &rootFS); err != nil {
		return err
	}
	rootFS.DiffIDs = append(rootFS.DiffIDs, newLayerDigest)
	p, err = json.Marshal(rootFS)
	if err != nil {
		return err
	}
	m["rootfs"] = p

	// Convert our entire image configuration back to bytes, and write it to the content store.
	p, err = json.Marshal(m)
	if err != nil {
		return err
	}
	manifest.Config.Digest = digest.FromBytes(p)
	err = content.WriteBlob(ctx, session.client.ContentStore(),
		ref,
		bytes.NewReader(p),
		int64(len(p)),
		manifest.Config.Digest,
	)
	if err != nil {
		return err
	}

	return err
}

func (session *Session) createImageFromSnapshot(ctx context.Context, img containerd.Image, activeSnapshotKey string, newReference string) error {
	ctx, done, err := session.client.WithLease(ctx) // Prevent garbage collection while we work.
	if err != nil {
		return err
	}
	defer done()

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

	// Generate a diff in content store
	diffs, err := session.client.DiffService().DiffMounts(ctx,
		lowerMounts,
		upperMounts,
		diff.WithMediaType(ocispec.MediaTypeImageLayerGzip),
		diff.WithReference("custom-ref"))
	if err != nil {
		return err
	}

	// Add our new layer to the image manifest
	manifest.Layers = append(manifest.Layers, diffs)

	// Add the blob checksum to image config
	info, err := contentStore.Info(ctx, diffs.Digest)
	if err != nil {
		return err
	}
	diffIDStr, ok := info.Labels[containerdUncompressed]
	if !ok {
		return fmt.Errorf("invalid differ response with no diffID")
	}
	diffIDDigest, err := digest.Parse(diffIDStr)
	if err != nil {
		return err
	}
	err = session.patchImageConfig(ctx, "custom-ref", &manifest, diffIDDigest)
	if err != nil {
		return err
	}

	// Prepare the labels that will tell the garbage collector
	// to NOT delete the content this manifest references.
	labels := map[string]string{
		"containerd.io/gc.ref.content.0": manifest.Config.Digest.String(),
	}
	for i, layer := range manifest.Layers {
		labels[fmt.Sprintf("containerd.io/gc.ref.content.%d", i+1)] = layer.Digest.String()
	}

	// Save our new image manifest, which now hows our new layer,
	// and a patched image config with a reference to the new layer.
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	manifestDigest := digest.FromBytes(manifestBytes)
	if err := content.WriteBlob(ctx,
		contentStore,
		"custom-ref",
		bytes.NewReader(manifestBytes),
		int64(len(manifestBytes)),
		manifestDigest,
		content.WithLabels(labels)); err != nil {
		return err
	}

	// Let's see if the image exists already, if so, let's delete it
	_, err = session.client.GetImage(ctx, newReference)
	if err == nil {
		session.client.ImageService().Delete(ctx, newReference, images.SynchronousDelete())
	}

	_, err = session.client.ImageService().Create(ctx,
		images.Image{
			Name: newReference,
			Target: ocispec.Descriptor{
				Digest:    manifestDigest,
				Size:      int64(len(manifestBytes)),
				MediaType: ocispec.MediaTypeImageManifest,
			},
		})
	if err != nil {
		return err
	}

	// This will create the required snapshot for the new layer,
	// which will allow us to run the image immediately.
	newImage, err := session.client.GetImage(ctx, newReference)
	if err != nil {
		return err
	}
	err = newImage.Unpack(ctx, containerd.DefaultSnapshotter)
	if err != nil {
		return err
	}

	return nil
}
