package push

import (
	"fmt"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "push",
		Usage: "Push an image to Docker Hub.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "imageName",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Pushing...")
			return nil
		},
	}
}
