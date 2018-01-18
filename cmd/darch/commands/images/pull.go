package images

import (
	ctx "context"
	"fmt"

	"github.com/containerd/containerd/reference"
	"github.com/pauldotknopf/darch/repository"
	"github.com/urfave/cli"
)

var pullCommand = cli.Command{
	Name:        "pull",
	Usage:       "pull an image from a remote registry",
	ArgsUsage:   "<image>",
	Description: "Pull and prepare an image for use in darch.",
	Action: func(context *cli.Context) error {
		var (
			image = context.Args().First()
		)

		refspec, err := reference.Parse(image)
		if err != nil {
			return reference.ErrInvalid
		}
		if refspec.Object == "" {
			return reference.ErrObjectRequired
		}

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		fmt.Printf("Pulling %s\n", refspec.String())

		img, err := repo.Pull(ctx.Background(), image)
		if err != nil {
			return err
		}

		fmt.Printf("Pulled %s\n", img.Target().Digest)

		return nil
	},
}
