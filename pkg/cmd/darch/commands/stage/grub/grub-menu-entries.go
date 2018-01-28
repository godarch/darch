package grub

import (
	"github.com/pauldotknopf/darch/pkg/cmd/darch/commands"
	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/urfave/cli"
	"os"
)

var grubMenuEntriesCommand = cli.Command{
	Name:        "menu-entries",
	Description: "outut a menu entries for each staged image",
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
			err = session.PrintGrubMenuEntry(stagedImage, os.Stdout)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
