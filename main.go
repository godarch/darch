package main

import (
	"fmt"
	"github.com/urfave/cli" // imports as package "cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."
	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Build an image.",
			Action: func(c *cli.Context) error {
				fmt.Println("Building...")
				return nil
			},
		},
		{
			Name:  "extract",
			Usage: "Extract an image to a temporary directory.",
			Action: func(c *cli.Context) error {
				fmt.Println("Extracting...")
				return nil
			},
		},
		{
			Name:  "stage",
			Usage: "Stage an image for booting.",
			Action: func(c *cli.Context) error {
				fmt.Println("Staging...")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
