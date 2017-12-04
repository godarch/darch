package inspect

import (
	"fmt"
	"log"
	"strings"

	"../../images"
	"../../utils"
	"github.com/urfave/cli"
)

func parentsCommand() cli.Command {
	return cli.Command{
		Name:      "parents",
		Usage:     "The parents (inherited images) of an image.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir, d",
				Usage: "Location of the images.",
				Value: ".",
			},
			cli.BoolFlag{
				Name: "excludeExternal",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := parents(c.Args().First(), c.String("imagesDir"), c.Bool("excludeExternal"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "inspect",
		Usage:     "Inspect an image.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir, d",
				Usage: "Location of the images.",
				Value: ".",
			},
		},
		Subcommands: []cli.Command{
			parentsCommand(),
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := inspect(c.Args().First(), c.String("imagesDir"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func parents(name string, imagesDir string, excludeExternal bool) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	imageDefinition, err := images.BuildDefinition(name, imagesDir)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	parents := make([]*images.ImageDefinition, 0)

	current := imageDefinition

	for current != nil {
		parents = append(parents, current)
		if strings.HasPrefix(current.Inherits, "external:") {
			current = nil
		} else {
			current, err = images.BuildDefinition(current.Inherits, imagesDir)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}
	}

	for _, child := range parents[1:] {
		log.Println(child.Name)
	}

	if !excludeExternal {
		externalImage := parents[len(parents)-1].Inherits
		log.Println(externalImage[len("external:"):len(externalImage)])

	}
	return nil
}

func inspect(name string, imagesDir string) error {
	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	imageDefinition, err := images.BuildDefinition(name, imagesDir)

	log.Println(imageDefinition.Name)

	return err
}
