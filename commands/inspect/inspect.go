package inspect

import (
	"fmt"
	"log"

	"../../images"
	"../../utils"
	"github.com/disiqueira/gotree"
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

func treeCommand() cli.Command {
	return cli.Command{
		Name:  "tree",
		Usage: "Display all images in a tree.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir, d",
				Usage: "Location of the images.",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			err := tree(c.String("imagesDir"))
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
			treeCommand(),
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
		if current.InheritsExternal {
			if !excludeExternal {
				log.Println(current.Inherits)
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

	imageDefinitions, err := images.BuildAllDefinitions(imagesDir)

	if err != nil {
		return err
	}

	current, ok := imageDefinitions[name]

	if !ok {
		return fmt.Errorf("Image %s doesn't exist", name)
	}

	for _, imageDefinition := range imageDefinitions {
		if imageDefinition.Inherits == current.Name {
			log.Println(imageDefinition.Name)
		}
	}

	return err
}

func buildTreeRecursively(parentDefinition images.ImageDefinition, imageDefinitions map[string]images.ImageDefinition) []gotree.GTStructure {
	children := make([]gotree.GTStructure, 0)

	for _, childImageDefinition := range imageDefinitions {
		if childImageDefinition.Inherits == parentDefinition.Name {
			var childNode gotree.GTStructure
			childNode.Name = childImageDefinition.Name

			for _, child := range buildTreeRecursively(childImageDefinition, imageDefinitions) {
				childNode.Items = append(childNode.Items, child)
			}
			children = append(children, childNode)
		}
	}

	return children
}

func tree(imagesDir string) error {
	if len(imagesDir) == 0 {
		return fmt.Errorf("Images directory is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	imageDefinitions, err := images.BuildAllDefinitions(imagesDir)

	if err != nil {
		return err
	}

	externalImages := make([]string, 0)

	for _, imageDefinition := range imageDefinitions {
		if imageDefinition.InheritsExternal {
			externalImages = append(externalImages, imageDefinition.Inherits)
		}
	}

	// this will be our root items
	externalImages = utils.RemoveDuplicates(externalImages)

	var rootNode gotree.GTStructure

	for _, externalImage := range externalImages {
		var externalImageNode gotree.GTStructure
		externalImageNode.Name = externalImage
		for _, imageDefinition := range imageDefinitions {
			if imageDefinition.InheritsExternal && imageDefinition.Inherits == externalImage {
				var childNode gotree.GTStructure
				childNode.Name = imageDefinition.Name
				for _, child := range buildTreeRecursively(imageDefinition, imageDefinitions) {
					childNode.Items = append(childNode.Items, child)
				}
				externalImageNode.Items = append(externalImageNode.Items, childNode)
			}
		}
		rootNode.Items = append(rootNode.Items, externalImageNode)
	}

	gotree.PrintTree(rootNode)

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
