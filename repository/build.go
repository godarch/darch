package repository

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/opencontainers/image-spec/identity"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/reference"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
	"github.com/pauldotknopf/darch/workspace"
)

// BuildRecipe Builds a recipe.
func (session *Session) BuildRecipe(ctx context.Context, recipe recipes.Recipe, tag string, buildPrefix string, environmentVariables map[string]string) error {

	ctx = namespaces.WithNamespace(ctx, "darch")

	if len(tag) == 0 {
		tag = "local"
	}

	session.client.ContentStore().Walk(ctx, func(content content.Info) error {
		fmt.Println(content.Digest.String())
		if content.Digest.String() == "sha256:1f0f5c30de52c731c9069d635337ccaa35f23042e81ebc25a77d13449cb9c19a" {
			fmt.Println("found")
			//session.client.DiffService().Apply()
			var f = "sdf"
			fmt.Println(f)
		}
		return nil
	})

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

	ds, err := img.RootFS(ctx)
	if err != nil {
		return err
	}
	for _, d := range ds {
		fmt.Println(d.String())
	}

	t := img.Target()
	fmt.Printf("Digest %", t.Digest.String())

	// im := images.Image{
	// 	Name:   "new-image:latest",
	// 	Target: img.Target(),
	// 	Labels: map[string]string{
	// 		"containerd.io/checkpoint": "true",
	// 	},
	// }

	// img2, err := session.client.ImageService().Create(ctx, im)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(img2.Name)

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

	err = session.client.SnapshotService(containerd.DefaultSnapshotter).Commit(ctx, "test-commit-name", snapshotKey)
	if err != nil {
		return err
	}

	// TODO: save image
	//

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
	return nil
}

func (session *Session) deleteSnapshot(ctx context.Context, snapshotKey string) error {
	return session.client.SnapshotService(containerd.DefaultSnapshotter).Remove(ctx, snapshotKey)
}

// func (session *Session) runRecipeBuild(ctx context.Context, recipe recipes.Recipe, snapshotKey string, img containerd.Image) error {
// 	id := utils.NewID()
// 	container, err := session.client.NewContainer(ctx,
// 		id,
// 		containerd.WithImage(img),
// 		containerd.WithSnapshotter(containerd.DefaultSnapshotter),
// 		containerd.WithSnapshot(snapshotKey),
// 		containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
// 		containerd.WithNewSpec(
// 			oci.WithImageConfig(img),
// 			oci.WithHostNamespace(specs.NetworkNamespace),
// 			oci.WithHostResolvconf,
// 			oci.WithMounts([]specs.Mount{
// 				specs.Mount{
// 					Destination: "/recipes",
// 					Type:        "bind",
// 					Source:      recipe.RecipesDir,
// 					Options:     []string{"rbind", "ro"},
// 				},
// 			}),
// 			oci.WithProcessArgs("/usr/bin/env", "bash", "-c", fmt.Sprintf("/darch-prepare && ./darch-runrecipe %s && ./darch-teardown", recipe.Name))))
// 	if err != nil {
// 		return err
// 	}

// 	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

// 	t, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
// 	if err != nil {
// 		return err
// 	}

// 	err = t.Start(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer t.Delete(ctx)

// 	var statusC <-chan containerd.ExitStatus
// 	if statusC, err = t.Wait(ctx); err != nil {
// 		return err
// 	}

// 	sigc := commands.ForwardAllSignals(ctx, t)
// 	defer commands.StopCatch(sigc)

// 	status := <-statusC
// 	code, _, err := status.Result()
// 	if err != nil {
// 		return err
// 	}

// 	if code != 0 {
// 		return cli.NewExitError("", int(code))
// 	}

// 	return err
// }
