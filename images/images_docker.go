package images

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// BuildImageLayer Run installation scripts on top of another image.
func BuildImageLayer(imageDefinition *ImageDefinition) error {

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	inherits := imageDefinition.Inherits[0]
	if strings.HasPrefix(inherits, "external:") {
		inherits = inherits[len("external:"):len(inherits)]
	}

	container, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: inherits,
	}, nil, nil, "building-"+imageDefinition.Name)

	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})

	if err != nil {
		destroyContainer(container.ID, cli)
		return err
	}

	_, err = cli.ContainerCommit(context.Background(),
		container.ID,
		types.ContainerCommitOptions{
			Reference: imageDefinition.Name,
		},
	)

	if err != nil {
		destroyContainer(container.ID, cli)
		return err
	}

	err = cli.ContainerStop(context.Background(), container.ID, nil)

	if err != nil {
		destroyContainer(container.ID, cli)
		return err
	}

	err = cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})

	if err != nil {
		destroyContainer(container.ID, cli)
		return err
	}

	return nil
}

func destroyContainer(containerId string, cli *client.Client) error {

	// TODO: Determine if container is running already.

	err := cli.ContainerStop(context.Background(), containerId, nil)

	if err != nil {
		log.Println("There was an error stopping the container.")
		return err
	}

	err = cli.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{})

	if err != nil {
		log.Println("There was an error removing the container.")
		return err
	}

	return nil
}
