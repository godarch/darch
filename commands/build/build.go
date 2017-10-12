package build

import (
	"fmt"
	"log"
	"strings"

	"../../images"
	"../../utils"
	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name: "build",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tags",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "imagePrefix",
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
			err := build(c.Args().First(), c.String("imagesDir"), strings.Split(c.String("tags"), ","), c.String("imagePrefix"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return err
		},
	}
}

func build(name string, imagesDir string, tags []string, imagePrefix string) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	log.Println("Name: " + name)
	log.Println("Images directory: " + imagesDir)
	log.Println("Image prefix: " + imagePrefix)

	if len(tags) > 0 {
		log.Println("Tags:")
		for _, tag := range tags {
			log.Println("\t" + tag)
		}
	} else {
		log.Println("Tags: none")
	}

	imageDefinition, err := images.BuildDefinition(name, imagesDir)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = images.BuildImageLayer(
		imageDefinition,
		tags,
		imagePrefix)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
