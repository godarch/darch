package stage

import (
	"fmt"

	"github.com/pauldotknopf/darch/pkg/staging"
	"github.com/urfave/cli"
)

var grubCommand = cli.Command{
	Name:   "grub",
	Usage:  "helper command for generating grub.cfd",
	Hidden: true,
	Action: func(clicontext *cli.Context) error {
		err := checkForRoot()
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
			fmt.Printf("%s %s %s %s %s %s\n", stagedImage.Ref.Name, stagedImage.Ref.Tag, stagedImage.Dir, stagedImage.Kernel, stagedImage.InitRAMFS, stagedImage.RootFS)
		}
		return nil
	},
}
