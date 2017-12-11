package builddep

import (
	"fmt"
	"log"

	"../../images"
	"../../utils"
	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "build-dep",
		Usage:     "List the dependencies needed to build images.",
		ArgsUsage: "IMAGE_NAME*N",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "images-dir, d",
				Usage: "Location of the images.",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			err := buildDep(c.Args(), c.String("images-dir"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return err
		},
	}
}

func walkImageRecursively(imageDefinition images.ImageDefinition, imageDefinitions map[string]images.ImageDefinition) []string {
	result := make([]string, 0)
	result = append(result, imageDefinition.Name)
	if !imageDefinition.InheritsExternal {
		children := walkImageRecursively(imageDefinitions[imageDefinition.Inherits], imageDefinitions)
		result = append(result, children...)
	}
	return result
}

func buildDep(imageNames []string, imagesDir string) error {

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

	dependencies := make([]string, 0)

	for _, imageDefinition := range imageDefinitions {
		if len(imageNames) == 0 || utils.Contains(imageNames, imageDefinition.Name) {
			parents := walkImageRecursively(imageDefinition, imageDefinitions)
			parents = utils.Reverse(parents)
			dependencies = append(dependencies, parents...)
		}
	}

	dependencies = utils.RemoveDuplicates(dependencies)

	for _, dependency := range dependencies {
		log.Println(dependency)
	}

	return nil
}
