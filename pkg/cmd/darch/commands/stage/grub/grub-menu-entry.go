package grub

import (
	"fmt"
	"github.com/pauldotknopf/darch/pkg/cmd/darch/commands"
	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/urfave/cli"
	"os"
)

var grubMenuEntryCommand = cli.Command{
	Name:        "menu-entry",
	Description: "outut a menu entry for a staged item",
	ArgsUsage:   "<image[:tag]>",
	Action: func(clicontext *cli.Context) error {
		var (
			imageName = clicontext.Args().First()
		)

		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		imageRef, err := reference.ParseImage(imageName)
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
			if stagedImage.Ref.FullName() == imageRef.FullName() {
				return session.PrintGrubMenuEntry(stagedImage, os.Stdout)
			}
		}

		return fmt.Errorf("image not found")
	},
}
