package stage

import (
	"fmt"

	"github.com/godarch/darch/pkg/cmd/darch/commands"
	"github.com/godarch/darch/pkg/staging"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list all staged images",
	Action: func(clicontext *cli.Context) error {
		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		session, err := staging.NewSession()
		if err != nil {
			return err
		}

		stagedImages, err := session.GetAllStaged()
		if err != nil {
			return err
		}

		for _, stagedImage := range stagedImages {
			fmt.Println(stagedImage.Ref.FullName())
		}
		return nil
	},
}
