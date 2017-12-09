package extract

import (
	"fmt"
	"log"
	"path"

	"../../images"
	"../../utils"
	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "extract",
		Usage:     "Extract an image.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tag, t",
				Usage: "The tag to extract.",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "destination, d",
				Usage: "The location to extract the image to.",
				Value: "/var/darch/extracted",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := extract(c.Args().First(), c.String("tag"), c.String("destination"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return err
		},
	}
}

func extract(name string, tag string, destinationDirectory string) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(tag) == 0 {
		return fmt.Errorf("Tag is required")
	}

	if len(destinationDirectory) == 0 {
		return fmt.Errorf("Destination is required")
	}

	destinationDirectory = utils.ExpandPath(destinationDirectory)
	destinationDirectory = path.Join(destinationDirectory, name+"/"+tag)

	log.Println("Name: " + name)
	log.Println("Tag: " + tag)
	log.Println("Destination: " + destinationDirectory)

	return images.ExtractImage(name, tag, destinationDirectory)
}
