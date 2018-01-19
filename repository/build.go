package repository

import (
	gocontext "context"
	"fmt"
	"runtime"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/oci"
	"github.com/urfave/cli"

	"github.com/containerd/containerd"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/reference"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
)

// BuildRecipe Builds a recipe.
func (session *Session) BuildRecipe(context gocontext.Context, recipe recipes.Recipe, tag string, buildPrefix string, environmentVariables map[string]string) error {

	context = namespaces.WithNamespace(context, "darch")

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

	img, err := session.client.GetImage(context, inheritsRef.String())
	if err != nil {
		// maybe it was because we don't have it? let's try to fetch it
		img, err = session.Pull(context, inheritsRef.String())
		if err != nil {
			return err
		}
	}

	id := utils.NewID()
	cntner, err := session.client.NewContainer(context,
		id,
		containerd.WithImage(img),
		containerd.WithSnapshotter(containerd.DefaultSnapshotter),
		containerd.WithNewSnapshot(id, img),
		containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), nil),
		containerd.WithNewSpec(
			oci.WithImageConfig(img),
			oci.WithProcessArgs("ls")))
	if err != nil {
		return err
	}

	defer cntner.Delete(context, containerd.WithSnapshotCleanup)

	t, err := cntner.NewTask(context, cio.Stdio)
	if err != nil {
		return err
	}

	err = t.Start(context)
	if err != nil {
		return err
	}
	defer t.Delete(context)

	var statusC <-chan containerd.ExitStatus
	if statusC, err = t.Wait(context); err != nil {
		return err
	}

	sigc := commands.ForwardAllSignals(context, t)
	defer commands.StopCatch(sigc)

	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}

	if _, err := t.Delete(context); err != nil {
		return err
	}
	if code != 0 {
		return cli.NewExitError("", int(code))
	}

	return err
}
