package repository

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
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

	tempMountsWs, err := workspace.NewWorkspace("")
	if err != nil {
		return err
	}
	defer tempMountsWs.Destroy()

	mounts, err := createTempMounts(tempMountsWs.Path)

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

	upperMounts, err := session.snapshotter.Mounts(ctx, snapshotKey)
	if err != nil {
		return err
	}

	err = mount.WithTempMount(ctx, upperMounts, func(root string) error {
		srcDir := path.Join(root, "extract")
		return filepath.Walk(path.Join(root, "extract"), func(_path string, _f os.FileInfo, _err error) error {
			if _err != nil {
				return _err
			}
			if !_f.IsDir() && strings.HasPrefix(_path, srcDir) {
				return utils.CopyFile(_path, path.Join(destination, _path[len(srcDir):]))
			}
			return nil
		})
	})
	if err != nil {
		return err
	}

	return nil
}
