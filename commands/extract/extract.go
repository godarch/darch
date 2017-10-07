package extract

import (
	"log"
	"path"

	"../../images"
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
				Name:  "imagesDir",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tag",
				Value: "local",
			},
			cli.StringFlag{
				Name: "destination",
			},
		},
		Action: func(c *cli.Context) error {
			return extract(c.Args().First(), c.String("imagesDir"), c.String("tag"), c.String("destination"))
		},
	}
}

func extract(imageName string, imagesDir string, tag string, destination string) error {
	log.Println("Images directory: " + imagesDir)
	log.Println("Image name: " + imageName)
	log.Println("Tag: " + tag)

	imageDefinition, err := images.BuildDefinition(imageName, imagesDir)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(destination) == 0 {
		destination = path.Join(imageDefinition.ImagesDir, ".extracted", imageDefinition.Name+"-"+tag)
	}

	log.Println("Destination: " + destination)

	images.ExtractImage(imageDefinition, tag, destination)

	return nil
}
