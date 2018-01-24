package stage

import (
	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/staging"
	"github.com/urfave/cli"
)

var runHooksCommand = cli.Command{
	Name:      "run-hooks",
	Usage:     "run hooks for image(s)",
	ArgsUsage: "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		var (
			imageName = clicontext.Args().First()
		)

		err := checkForRoot()
		if err != nil {
			return err
		}

		stagingSession, err := staging.NewSession()
		if err != nil {
			return err
		}

		if len(imageName) > 0 {
			imageRef, err := reference.ParseImage(imageName)
			if err != nil {
				return err
			}
			return stagingSession.RunHooksForImage(imageRef)
		}

		return stagingSession.RunAllHooks()
	},
}
