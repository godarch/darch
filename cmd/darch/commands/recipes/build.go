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
			recipeName = clicontext.Args().First()
		)

		fmt.Printf("Building %s...\n", recipeName)

		return nil
	},
}
