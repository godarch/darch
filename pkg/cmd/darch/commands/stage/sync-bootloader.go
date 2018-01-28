package stage

import (
	"github.com/pauldotknopf/darch/pkg/cmd/darch/commands"
	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/urfave/cli"
)

var syncBootloaderCommand = cli.Command{
	Name:      "sync-bootloader",
	Usage:     "Updates",
	ArgsUsage: "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		stagingSession, err := staging.NewSession()
		if err != nil {
			return err
		}

		return stagingSession.SyncBootloader()
	},
}
