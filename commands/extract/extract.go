package extract

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "extract",
		Usage:     "Extract an image.",
		UsageText: "darch extract [options] [image]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imageDir",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tag",
				Value: "latest",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Extracting...")
			return nil
		},
	}
}

func extract(imageName string, imageDir string, flags []string) error {
	log.Println("Image directory: " + imageDir)
	log.Println("Image name: " + imageName)

	// imageDefinition, err := images.BuildDefinition(imageName, imageDir)

	// if err != nil {
	// 	return cli.NewExitError(err, 1)
	// }

	return nil
}
