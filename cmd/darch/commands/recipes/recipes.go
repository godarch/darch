package recipes

import (
	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:  "recipes",
		Usage: "view/build recipes",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "recipes-dir, d",
				Usage: "location of the recipes",
				Value: ".",
			},
		},
		Subcommands: cli.Commands{
			buildCommand,
			listCommand,
			parentsCommand,
		},
	}
)

func getRecipesDir(ctx *cli.Context) string {
	return ctx.GlobalString("recipes-dir")
}
