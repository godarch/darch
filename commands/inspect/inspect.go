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

func childrenCommand() cli.Command {
	return cli.Command{
		Name:      "children",
		Usage:     "The children that are dependent on the provided image.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir, d",
				Usage: "Location of the images.",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := children(c.Args().First(), c.String("imagesDir"))
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
			childrenCommand(),
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

	imageDefinitions, err := images.BuildAllDefinitions(imagesDir)

	if err != nil {
		return err
	}

	current, ok := imageDefinitions[name]

	if !ok {
		return fmt.Errorf("Image %s doesn't exist", name)
	}

	finished := false
	for finished != true {
		if strings.HasPrefix(current.Inherits, "external:") {
			if !excludeExternal {
				log.Println(current.Inherits[len("external:"):len(current.Inherits)])
			}
			finished = true
		} else {
			current = imageDefinitions[current.Inherits]
			log.Println(current.Name)
		}
	}

	return nil
}

func children(name string, imagesDir string) error {
	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	allDefinitions, err := images.BuildAllDefinitions(imagesDir)

	for _, imageDefinition := range allDefinitions {
		log.Println(imageDefinition.Name)
	}

	return err
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
