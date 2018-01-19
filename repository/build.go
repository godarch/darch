package repository

import (
	"context"
	"fmt"
	"runtime"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/image-spec/identity"
	"github.com/urfave/cli"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/reference"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
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

	// Let's create the snapshot that all of our containers will run off of
	snapshotKey := utils.NewID()
	session.createSnapshot(ctx, snapshotKey, img)
	defer session.deleteSnapshot(ctx, snapshotKey)

	// Testing, to see if we can run multiple containers against the same snapshot.
	if err = session.runContainer(ctx, snapshotKey, img, "touch", "/test"); err != nil {
		return err
	}

	if err = session.runContainer(ctx, snapshotKey, img, "ls /"); err != nil {
		return err
	}

	return err
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

	mounts, err := session.client.SnapshotService(containerd.DefaultSnapshotter).Mounts(ctx, snapshotKey)

	for _, m := range mounts {
		fmt.Println(m.Source)
		fmt.Println(m.Type)
		fmt.Println(m.Options)
	}

	return nil
}

func (session *Session) runContainer(ctx context.Context, snapshotKey string, img containerd.Image, args ...string) error {
	id := utils.NewID()
	container, err := session.client.NewContainer(ctx,
		id,
		containerd.WithImage(img),
		containerd.WithSnapshotter(containerd.DefaultSnapshotter),
		containerd.WithSnapshot(snapshotKey),
		containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
		containerd.WithNewSpec(
			oci.WithImageConfig(img),
			oci.WithProcessArgs(args...)))
	if err != nil {
		return err
	}

	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	t, err := container.NewTask(ctx, cio.Stdio)
	if err != nil {
		return err
	}

	err = t.Start(ctx)
	if err != nil {
		return err
	}
	defer t.Delete(ctx)

	var statusC <-chan containerd.ExitStatus
	if statusC, err = t.Wait(ctx); err != nil {
		return err
	}

	sigc := commands.ForwardAllSignals(ctx, t)
	defer commands.StopCatch(sigc)

	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}

	if code != 0 {
		return cli.NewExitError("", int(code))
	}

	return err
}

func (session *Session) deleteSnapshot(ctx context.Context, snapshotKey string) error {
	return session.client.SnapshotService(containerd.DefaultSnapshotter).Remove(ctx, snapshotKey)
}
