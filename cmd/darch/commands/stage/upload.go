package stage

import (
	"context"

	"github.com/pauldotknopf/darch/repository"
	"github.com/pauldotknopf/darch/workspace"
	"github.com/urfave/cli"
)

var uploadCommand = cli.Command{
	Name:      "upload",
	Usage:     "upload local image to stage",
	ArgsUsage: "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		var (
			imageName = clicontext.Args().First()
		)

		err := checkForRoot()
		if err != nil {
			return err
		}

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		ws, err := workspace.NewWorkspace("/var/lib/darch/tmp")
		if err != nil {
			return err
		}
		defer ws.Destroy()

		err = repo.ExtractImage(context.Background(), imageName, ws.Path)
		if err != nil {
			return err
		}

		return nil
	},
}
