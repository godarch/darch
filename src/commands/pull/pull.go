package pull

import (
	"fmt"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "pull",
		Usage: "Pull an image from Docker Hub.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "imageName",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Pulling...")
			return nil
		},
	}
}
