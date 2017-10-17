package stage

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"../../utils"
	"github.com/urfave/cli"
)

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:      "stage",
		Usage:     "Stage an image for booting.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tag, t",
				Usage: "The tag to stage.",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "source, s",
				Usage: "The location where the extract images are.",
				Value: "/var/darch",
			},
			cli.StringFlag{
				Name:  "fstab",
				Usage: "The fstab file to use for the booted image. If relative path, darch will look in \"source\" for the file.",
				Value: "defaultfstab",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := stage(c.Args().First(), c.String("tag"), c.String("source"), c.String("fstab"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func stage(name string, tag string, sourceDirectory string, fstab string) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(tag) == 0 {
		return fmt.Errorf("Tag is required")
	}

	if len(sourceDirectory) == 0 {
		return fmt.Errorf("Source is required")
	}

	if len(fstab) == 0 {
		return fmt.Errorf("fstab is required")
	}

	sourceDirectory = utils.ExpandPath(sourceDirectory)

	log.Println("Name: " + name)
	log.Println("Tag: " + tag)
	log.Println("Source: " + sourceDirectory)

	sourceImageDirectory := path.Join(sourceDirectory, name+"/"+tag)

	if !utils.DirectoryExists(sourceImageDirectory) {
		return fmt.Errorf("No image found at %s", sourceImageDirectory)
	}

	destinationDirectory := path.Join("/boot", "darch", name, tag)

	if utils.DirectoryExists(destinationDirectory) {
		log.Println("Cleaning already existing staging directory...")
		err := os.RemoveAll(destinationDirectory)
		if err != nil {
			log.Printf("Couldn't delete already staged directory %s\n", destinationDirectory)
			return err
		}
	}

	log.Printf("Copying boot files to %s...\n", destinationDirectory)

	err := utils.CopyDir(sourceImageDirectory, destinationDirectory)
	if err != nil {
		return err
	}

	if len(fstab) > 0 {
		if !path.IsAbs(fstab) {
			fstab = utils.ExpandPath(path.Join(sourceDirectory, fstab))
		}

		log.Printf("Staging fstab file %s\n", fstab)

		err = utils.CopyFile(fstab, path.Join(destinationDirectory, "fstab"))

		if err != nil {
			return err
		}
	}

	log.Printf("Successfully staged %s at tag %s\n", name, tag)

	if !utils.FileExists("/etc/grub.d/60_darch") {
		return fmt.Errorf("Grub generator doesn't exist at %s", "/etc/grub.d/60_darch")
	}

	log.Println("Generating grub boot entries...")

	cmd := exec.Command("grub-mkconfig", "--output", "/boot/grub/grub.cfg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err == nil {
		return err
	}

	log.Println("Finished!")

	return nil
}
