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
		ArgsUsage: "IMAGE_NAME*N",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "images-dir, d",
				Usage: "Location of the images to build.",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tags, t",
				Usage: "Command separated list of tags to associated with the built image.",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "image-prefix, p",
				Usage: "Prefix for built images. For example, a value of \"pauldotknopf/darch-\" while building image \"base\", the generated image will be named \"pauldotknopf/darch-base\".",
				Value: "",
			},
			cli.StringFlag{
				Name:  "package-cache, c",
				Usage: "Location where package caches are stored. This speeds up builds, preventing downloading.",
				Value: "/var/darch/cache/packages",
			},
			cli.StringSliceFlag{
				Name: "environment, e",
			},
		},
		Action: func(c *cli.Context) error {
			environmentVaribles, err := utils.ConvertVariableStringsToMap(c.StringSlice("environment"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			err = build(c.Args(), c.String("images-dir"), strings.Split(c.String("tags"), ","), c.String("image-prefix"), c.String("package-cache"), environmentVaribles)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return err
		},
	}
}

func build(imageNames []string, imagesDir string, tags []string, imagePrefix string, packageCache string, environmentVariables map[string]string) error {
	if len(imageNames) == 0 {
		return fmt.Errorf("You must provide at least one image")
	}

	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	imageDefinitions, err := images.BuildAllDefinitions(imagesDir)
	if err != nil {
		return err
	}

	// Make sure the provided images exist
	for _, imageName := range imageNames {
		_, ok := imageDefinitions[imageName]
		if !ok {
			return fmt.Errorf("Image %s doesn't exist", imageName)
		}
	}

	log.Println("Name: " + strings.Join(imageNames, ","))
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

	for _, imageName := range imageNames {
		err = images.BuildImageLayer(
			imageDefinitions[imageName],
			tags,
			imagePrefix,
			packageCache,
			environmentVariables)

		if err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}
