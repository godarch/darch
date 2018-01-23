package repository

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"runtime"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/utils"
	"github.com/pauldotknopf/darch/workspace"
)

// ExtractImage Extracts an image (with tag) to a specified directory
func (session *Session) ExtractImage(ctx context.Context, name string, destination string) error {
	ctx = namespaces.WithNamespace(ctx, "darch")

	ctx, done, err := session.client.WithLease(ctx) // Prevent garbage collection while we work.
	if err != nil {
		return err
	}
	defer done()

	imgRef, err := reference.ParseImage(name)
	if err != nil {
		return err
	}

	img, err := session.client.GetImage(ctx, imgRef.FullName())
	if err != nil {
		return err
	}

	ws, err := workspace.NewWorkspace("/tmp")
	if err != nil {
		return err
	}
	defer ws.Destroy()

	mounts, err := createTempMounts(ws.Path)

	mounts = append(mounts, specs.Mount{
		Source:      "/home/pknopf/git/darch/src/github.com/pauldotknopf/darch/rootfs/helpers/darch-extract",
		Destination: "/darch-extract",
		Type:        "bind",
		Options:     []string{"rbind", "rw"},
	})

	// Create the snapshot that our extraction will happen on.
	snapshotKey := utils.NewID()
	err = session.createSnapshot(ctx, snapshotKey, img)
	if err != nil {
		return err
	}
	defer session.deleteSnapshot(ctx, snapshotKey)

	err = session.RunContainer(ctx, ContainerConfig{
		newOpts: []containerd.NewContainerOpts{
			containerd.WithImage(img),
			containerd.WithSnapshotter(containerd.DefaultSnapshotter),
			containerd.WithSnapshot(snapshotKey),
			containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
			containerd.WithNewSpec(
				oci.WithImageConfig(img),
				oci.WithHostNamespace(specs.NetworkNamespace),
				oci.WithMounts(mounts),
				oci.WithProcessArgs("/usr/bin/env", "bash", "-c", "/darch-extract"),
			),
		},
	})
	if err != nil {
		return err
	}

	snapshotService := session.client.SnapshotService(containerd.DefaultSnapshotter)
	diffService := session.client.DiffService()

	upperMounts, err := snapshotService.Mounts(ctx, snapshotKey)
	if err != nil {
		return err
	}

	desc, err := diffService.DiffMounts(ctx,
		[]mount.Mount{},
		upperMounts,
		diff.WithMediaType(ocispec.MediaTypeImageLayer),
		diff.WithReference("custom-ref"))

	fmt.Println(desc.Digest)

	rdrAt, err := session.client.ContentStore().ReaderAt(ctx, desc.Digest)
	if err != nil {
		return err
	}
	defer rdrAt.Close()

	tr := tar.NewReader(content.NewReader(rdrAt))

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		fmt.Println()
	}

	return nil
}
