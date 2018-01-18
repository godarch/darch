package recipes

import (
	"fmt"

	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:      "build",
	Usage:     "build a recipe",
	ArgsUsage: "<recipe>",
	Action: func(clicontext *cli.Context) error {
		var (
			recipe = clicontext.Args().First()
		)

		fmt.Printf("Building %s...\n", recipe)

		return nil
	},
}
