package images

import (
	"context"
	"fmt"

	"github.com/pauldotknopf/darch/repository"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list images",
	Action: func(clicontext *cli.Context) error {

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		imgs, err := repo.GetImages(context.Background())
		if err != nil {
			return err
		}

		for _, img := range imgs {
			fmt.Println(img.Name + ":" + img.Tag)
		}

		return nil
	},
}
