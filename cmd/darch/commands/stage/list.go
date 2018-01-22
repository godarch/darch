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
		stagedItems, err := staging.GetAllStaged(staging.DefaultStagingDirectory)
		if err != nil {
			return err
		}
		for _, stagedItem := range stagedItems {
			for _, stagedItemTag := range stagedItem.Tags {
				fmt.Println(stagedItem.Name + ":" + stagedItemTag.Name)
			}
		}
		return nil
	},
}
