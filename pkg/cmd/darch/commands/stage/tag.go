package stage

import (
	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/urfave/cli"
)

var tagCommand = cli.Command{
	Name:      "tag",
	Usage:     "tag images",
	ArgsUsage: "<src[:tag]> <dest[:tag]>",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "if overwriting existing tag, delete it",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			sourceImage      = clicontext.Args().First()
			destinationImage = clicontext.Args().Get(1)
			force            = clicontext.Bool("force")
		)

		sourceImageRef, err := reference.ParseImage(sourceImage)
		if err != nil {
			return err
		}

		destinationImageRef, err := reference.ParseImage(destinationImage)
		if err != nil {
			return err
		}

		stagingSession, err := staging.NewSession()
		if err != nil {
			return err
		}

		err = stagingSession.Tag(sourceImageRef, destinationImageRef, force)
		if err != nil {
			return err
		}

		// Run hooks for the newly tagged image.
		err = stagingSession.RunHooksForImage(destinationImageRef)
		if err != nil {
			return err
		}

		return stagingSession.SyncBootloader()
	},
}
