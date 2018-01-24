package stage

import (
	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/staging"
	"github.com/urfave/cli"
)

var removeCommand = cli.Command{
	Name:      "remove",
	Usage:     "removes an image from the stage",
	ArgsUsage: "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		var (
			imageName = clicontext.Args().First()
		)

		err := checkForRoot()
		if err != nil {
			return err
		}

		imageRef, err := reference.ParseImage(imageName)
		if err != nil {
			return err
		}

		stagingSession, err := staging.NewSession()
		if err != nil {
			return err
		}

		return stagingSession.Remove(imageRef)
	},
}
