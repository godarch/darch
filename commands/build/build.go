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
		Name:      "build",
		Usage:     "Build an image.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir, d",
				Usage: "Location of the images to build.",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tags, t",
				Usage: "Command separated list of tags to associated with the built image.",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "imagePrefix, p",
				Usage: "Prefix for built images. For example, a value of \"pauldotknopf/darch-\" while building image \"base\", the generated image will be named \"pauldotknopf/darch-base\".",
				Value: "",
			},
			cli.StringSliceFlag{
				Name: "environment, e",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			environmentVaribles, err := utils.ConvertVariableStringsToMap(c.StringSlice("environment"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			err = build(c.Args().First(), c.String("imagesDir"), strings.Split(c.String("tags"), ","), c.String("imagePrefix"), environmentVaribles)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return err
		},
	}
}

func build(name string, imagesDir string, tags []string, imagePrefix string, environmentVariables map[string]string) error {
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
		imagePrefix,
		environmentVariables)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
