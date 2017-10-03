package build

import (
	"log"

	"../../images"
	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "build",
		Usage:     "Build an image.",
		UsageText: "darch build [options] [image]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imageDir",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			return build(c.Args().First(), c.String("imageDir"))
		},
	}
}

func build(imageName string, imageDir string) error {
	log.Println("Image directory: " + imageDir)
	log.Println("Image name: " + imageName)

	imageDefinition, err := images.BuildDefinition(imageName, imageDir)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = images.BuildImageLayer(imageDefinition)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
