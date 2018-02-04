package stage

import (
	"github.com/godarch/darch/pkg/cmd/darch/commands"
	"github.com/godarch/darch/pkg/staging"
	"github.com/urfave/cli"
)

var cleanCommand = cli.Command{
	Name:  "clean",
	Usage: "clean any staged images that aren't referenced or booted",
	Action: func(clicontext *cli.Context) error {
		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		session, err := staging.NewSession()
		if err != nil {
			return err
		}

		return session.Clean()
	},
}
