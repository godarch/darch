package stage

import (
	"fmt"

	"github.com/pauldotknopf/darch/staging"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list all staged images",
	Action: func(clicontext *cli.Context) error {
		err := checkForRoot()
		if err != nil {
			return err
		}
		stagedImages, err := staging.GetAllStaged()
		if err != nil {
			return err
		}
		for _, stagedImage := range stagedImages {
			fmt.Println(stagedImage.Name + ":" + stagedImage.Tag)
		}
		return nil
	},
}
