package stage

import (
	"fmt"
	"log"
	"path"

	"../../images"
	"../../stage"
	"github.com/kennygrant/sanitize"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

func uploadCommand() cli.Command {
	return cli.Command{
		Name:      "upload",
		Usage:     "Upload an image to the stage to be booted.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tag, t",
				Usage: "The tag to stage.",
				Value: "local",
			},
		},
		Action: func(c *cli.Context) error {
			err := upload(c.Args().First(), c.String("tag"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func listCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List the images currently staged.",
		Action: func(c *cli.Context) error {
			err := list()
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
		Name:  "stage",
		Usage: "Commands that help manage the stage.",
		Subcommands: []cli.Command{
			uploadCommand(),
			listCommand(),
		},
	}
}

func upload(name string, tag string) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(tag) == 0 {
		return fmt.Errorf("Tag is required")
	}

	destinationDirectory := "/var/darch/staged"
	destinationDirectory = path.Join(destinationDirectory, sanitize.Path(name+"/"+tag))

	log.Println("Name: " + name)
	log.Println("Tag: " + tag)
	log.Println("Destination: " + destinationDirectory)

	err := images.ExtractImage(name, tag, destinationDirectory)

	if err != nil {
		return err
	}

	return nil
}

func list() error {
	stagedItems, err := stage.GetAllStaged("/var/darch/staged")

	if err != nil {
		return err
	}

	result := []string{
		"Name | Tag | Path | Kernel | InitramFS | RootFS",
	}

	for _, stagedItem := range stagedItems {
		for _, stagedItemTag := range stagedItem.Tags {
			result = append(result, stagedItem.Name+" | "+stagedItemTag.Name+" | "+stagedItemTag.Path+" | "+stagedItemTag.BootKernel+" | "+stagedItemTag.BootInitRAMFS+" | "+stagedItemTag.BootRootFS)
		}
	}

	fmt.Println(columnize.SimpleFormat(result))

	return nil
}
