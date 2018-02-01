package images

import (
	ctx "context"
	"fmt"

	"github.com/godarch/darch/pkg/cmd/darch/commands"
	"github.com/godarch/darch/pkg/reference"
	"github.com/godarch/darch/pkg/repository"
	"github.com/urfave/cli"
)

var pullCommand = cli.Command{
	Name:      "pull",
	Usage:     "pull an image from a remote registry",
	ArgsUsage: "[flags] <image>",
	Flags:     commands.RegistryFlags,
	Action: func(clicontext *cli.Context) error {
		var (
			image = clicontext.Args().First()
		)

		imageRef, err := reference.ParseImage(image)
		if err != nil {
			return err
		}

		resolver, err := commands.GetResolver(clicontext)
		if err != nil {
			return err
		}

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		fmt.Printf("pulling %s\n", imageRef.FullName())

		err = repo.Pull(ctx.Background(), imageRef, resolver)
		if err != nil {
			return err
		}

		fmt.Printf("pulled %s\n", imageRef.FullName())

		return nil
	},
}
