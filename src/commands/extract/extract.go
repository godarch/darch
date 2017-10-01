package extract

import (
	"fmt"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "extract",
		Usage: "Extract an image.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "imageName",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Extracting...")
			return nil
		},
	}
}
