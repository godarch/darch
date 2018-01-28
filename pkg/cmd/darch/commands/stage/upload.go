package stage

import (
	"context"

	"github.com/pauldotknopf/darch/pkg/cmd/darch/commands"
	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/repository"
	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/pauldotknopf/darch/pkg/workspace"
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

		err := commands.CheckForRoot()
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

		stagingSession, err := staging.NewSession()
		if err != nil {
			return err
		}

		ws, err := workspace.NewWorkspace(staging.DefaultStagingDirectoryTmp)
		if err != nil {
			return err
		}
		defer ws.Destroy()

		err = repo.ExtractImage(context.Background(), imageRef, ws.Path)
		if err != nil {
			return err
		}

		err = stagingSession.UploadDirectoryWithMove(ws.Path, imageRef, force)
		if err != nil {
			return err
		}
		ws.MarkDestroyed() // prevent defered Destroy() from working, since we moved the directory to where it should be.

		// Run hooks for the new image.
		return stagingSession.RunHooksForImage(imageRef)
	},
}
