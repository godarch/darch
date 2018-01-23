package stage

import (
	"context"
	"fmt"

	"github.com/pauldotknopf/darch/repository"
	"github.com/pauldotknopf/darch/staging"
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

		ws, err := workspace.NewWorkspace(staging.DefaultStagingDirectoryTmp)
		if err != nil {
			return err
		}
		defer ws.Destroy()

		fmt.Printf("Exporting image to %s\n", ws.Path)

		err = repo.ExtractImage(context.Background(), imageName, ws.Path)
		if err != nil {
			return err
		}

		err = staging.UploadDirectoryWithMove(ws.Path, imageName)
		if err != nil {
			return err
		}
		ws.MarkDestroyed() // prevent defered Destroy() from working, since we moved the directory to where it should be.

		return nil
	},
}
