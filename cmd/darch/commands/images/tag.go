package images

import (
	"context"
	"fmt"

	"github.com/pauldotknopf/darch/repository"
	"github.com/urfave/cli"
)

var tagCommand = cli.Command{
	Name:      "tag",
	Usage:     "tag images",
	ArgsUsage: "<src[:tag]> <dest[:tag]>",
	Action: func(clicontext *cli.Context) error {
		if len(clicontext.Args()) != 2 {
			return fmt.Errorf("invalid args")
		}
		var (
			sourceImage      = clicontext.Args().First()
			destinationImage = clicontext.Args().Get(1)
		)
		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		err = repo.TagImage(context.Background(), sourceImage, destinationImage)
		if err != nil {
			return err
		}

		return nil
	},
}
