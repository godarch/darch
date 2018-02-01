package images

import (
	"context"

	"github.com/godarch/darch/pkg/reference"
	"github.com/godarch/darch/pkg/repository"
	"github.com/urfave/cli"
)

var removeCommand = cli.Command{
	Name:      "remove",
	Usage:     "remove an image",
	ArgsUsage: "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		var (
			image = clicontext.Args().First()
		)

		ref, err := reference.ParseImage(image)
		if err != nil {
			return nil
		}

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		err = repo.RemoveImage(context.Background(), ref.FullName())
		if err != nil {
			return err
		}
		return nil
	},
}
