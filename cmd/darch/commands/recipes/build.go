package recipes

import (
	"fmt"

	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:        "build",
	Usage:       "build a recipe",
	ArgsUsage:   "<recipe>",
	Description: "Build a recipe.",
	Action: func(context *cli.Context) error {
		var (
			recipe = context.Args().First()
		)

		fmt.Printf("Building %s...\n", recipe)

		return nil
	},
}
