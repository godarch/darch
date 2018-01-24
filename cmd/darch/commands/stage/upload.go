package stage

import (
	"context"

	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/repository"
	"github.com/pauldotknopf/darch/staging"
	"github.com/pauldotknopf/darch/workspace"
	"github.com/urfave/cli"
)

var uploadCommand = cli.Command{
	Name:      "upload",
	Usage:     "upload local image to stage",
	ArgsUsage: "<image[:tag]>",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "overwrite existing image with the given name",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			imageName = clicontext.Args().First()
			force     = clicontext.Bool("force")
		)

		err := checkForRoot()
		if err != nil {
			return err
		}

		imageRef, err := reference.ParseImage(imageName)
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

		err = repo.ExtractImage(context.Background(), imageRef, ws.Path)
		if err != nil {
			return err
		}

		err = staging.UploadDirectoryWithMove(ws.Path, imageRef, force)
		if err != nil {
			return err
		}
		ws.MarkDestroyed() // prevent defered Destroy() from working, since we moved the directory to where it should be.

		return nil
	},
}
