package repository

import (
	"context"
	"path"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pauldotknopf/darch/utils"
	"github.com/urfave/cli"
)

// ContainerConfig configuration about how to run the container
type ContainerConfig struct {
	env     []string
	newOpts []containerd.NewContainerOpts
	delOpts []containerd.DeleteOpts
}

func createTempMounts(dir string) ([]specs.Mount, error) {

	mounts := []specs.Mount{}

	if utils.FileExists("/etc/resolv.conf") {
		err := utils.CopyFile("/etc/resolv.conf", path.Join(dir, "resolv.conf"))
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, specs.Mount{
			Destination: "/etc/resolv.conf",
			Type:        "bind",
			Source:      path.Join(dir, "resolv.conf"),
			Options:     []string{"rbind", "rw"},
		})
	}

	return mounts, nil
}

// RunContainer Runs a container
func (session *Session) RunContainer(ctx context.Context, config ContainerConfig) error {
	ctx = namespaces.WithNamespace(ctx, "darch")
	id := utils.NewID()
	container, err := session.client.NewContainer(ctx,
		id,
		config.newOpts...,
	)
	if err != nil {
		return err
	}

	defer container.Delete(ctx, config.delOpts...)

	t, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
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
		return cli.NewExitError("Error running container", int(code))
	}

	return err
}
