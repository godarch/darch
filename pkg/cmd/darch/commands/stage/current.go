package stage

import (
	"fmt"

	"github.com/godarch/darch/pkg/cmd/darch/commands"
	"github.com/godarch/darch/pkg/reference"
	"github.com/godarch/darch/pkg/staging"
	"github.com/urfave/cli"
)

var currentCommand = cli.Command{
	Name:  "current",
	Usage: "prints the current booted image",
	Action: func(clicontext *cli.Context) error {
		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		session, err := staging.NewSession()
		if err != nil {
			return err
		}

		current, err := session.GetCurrentBootedImage()

		if err == reference.ErrDoesNotExist {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println(current.Ref.FullName())

		return nil
	},
}
