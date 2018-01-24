package hooks

import (
	"fmt"

	"github.com/pauldotknopf/darch/hooks"
	"github.com/pauldotknopf/darch/staging"
	"github.com/urfave/cli"
)

var detailsCommand = cli.Command{
	Name:      "details",
	Usage:     "details about a hook configuration",
	ArgsUsage: "<hook>",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "include-matched-images",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			hookName             = clicontext.Args().First()
			includeMatchedImages = clicontext.Bool("include-matched-images")
		)

		hook, err := hooks.GetHook(hookName)
		if err != nil {
			return err
		}

		fmt.Printf("name: %s\n", hook.Name)
		fmt.Printf("path: %s\n", hook.Path)
		fmt.Printf("executionOrder: %d\n", hook.ExecutionOrder)
		fmt.Printf("include images:\n")
		for _, includeImage := range hook.IncludeImages {
			fmt.Printf("\t%s\n", includeImage)
		}
		fmt.Printf("exclude images:\n")
		for _, excludeImage := range hook.ExcludeImages {
			fmt.Printf("\t%s\n", excludeImage)
		}

		if includeMatchedImages {
			stagingSession, err := staging.NewSession()
			if err != nil {
				return err
			}

			fmt.Printf("matched images:\n")
			images, err := stagingSession.GetAllStaged()
			if err != nil {
				return err
			}
			for _, image := range images {
				if hooks.AppliesToImage(hook, image.Ref) {
					fmt.Printf("\t%s\n", image.Ref.FullName())
				}
			}
		}

		return nil
	},
}
