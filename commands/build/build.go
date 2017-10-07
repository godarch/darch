package build

import (
	"log"
	"strings"

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
				Name:  "imagesDir",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tags",
				Value: "local",
			},
		},
		Action: func(c *cli.Context) error {
			return build(c.Args().First(), c.String("imagesDir"), strings.Split(c.String("tags"), ","))
		},
	}
}

func build(imageName string, imagesDir string, flags []string) error {
	log.Println("Images directory: " + imagesDir)
	log.Println("Image name: " + imageName)

	imageDefinition, err := images.BuildDefinition(imageName, imagesDir)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = images.BuildImageLayer(imageDefinition,
		flags)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
