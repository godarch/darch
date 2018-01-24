package main

import (
	"fmt"
	"os"

	"github.com/pauldotknopf/darch/cmd/darch/commands/helpers"
	"github.com/pauldotknopf/darch/cmd/darch/commands/images"
	"github.com/pauldotknopf/darch/cmd/darch/commands/recipes"
	"github.com/pauldotknopf/darch/cmd/darch/commands/stage"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."
	app.Version = Version
	app.HideVersion = true
	app.Commands = []cli.Command{
		images.Command,
		recipes.Command,
		stage.Command,
		helpers.Command,
		cli.Command{
			Name:  "version",
			Usage: "Print version information about darch.",
			Action: func(c *cli.Context) error {
				fmt.Printf("version %s\n", Version)
				fmt.Printf("commit %s\n", GitCommit)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "darch: %s\n", err)
		os.Exit(1)
	}
}
