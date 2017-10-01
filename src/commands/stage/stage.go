package stage

import (
	"fmt"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "stage",
		Usage: "Stage an image for booting.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "imageName",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Staging...")
			return nil
		},
	}
}
