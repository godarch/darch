package build

import (
	"fmt"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "Build an image.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "imageName",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Building...")
			return nil
		},
	}
}
