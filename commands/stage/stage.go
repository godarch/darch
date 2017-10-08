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
		UsageText: "darch stage [options] [image]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "imagesDir",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "tag",
				Value: "local",
			},
			cli.StringFlag{
				Name:  "fstab",
				Value: "defaultfstab",
			},
		},
		Action: func(c *cli.Context) error {
			err := stage(c.Args().First(), c.String("imagesDir"), c.String("tag"), c.String("fstab"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func stage(imageName string, imagesDir string, tag string, fstab string) error {

	if len(imageName) == 0 {
		return fmt.Errorf("An image name is required")
	}

	imagesDir = utils.ExpandPath(imagesDir)

	if len(imagesDir) == 0 {
		return fmt.Errorf("An images directory must be provided")
	}

	log.Println("Image name: " + imageName)
	log.Println("Images directory: " + imagesDir)

	sourceDirectory := path.Join(imagesDir, ".extracted", imageName+"-"+tag)

	if !utils.DirectoryExists(sourceDirectory) {
		return fmt.Errorf("No image found at %s", sourceDirectory)
	}

	destinationDirectory := path.Join("/boot", "darch", imageName+"-"+tag)

	if utils.DirectoryExists(destinationDirectory) {
		log.Println("Cleaning already existing staging directory...")
		err := os.RemoveAll(destinationDirectory)
		if err != nil {
			log.Printf("Couldn't delete already staged directory %s\n", destinationDirectory)
			return err
		}
	}

	log.Printf("Copying boot files to %s...\n", destinationDirectory)

	err := utils.CopyDir(sourceDirectory, destinationDirectory)
	if err != nil {
		return err
	}

	if len(fstab) > 0 {
		if !path.IsAbs(fstab) {
			fstab = utils.ExpandPath(path.Join(imagesDir, fstab))
		}

		log.Printf("Staging fstab file %s\n", fstab)

		err = utils.CopyFile(fstab, path.Join(destinationDirectory, "fstab"))

		if err != nil {
			return err
		}
	}

	log.Printf("Successfully staged %s at tag %s\n", imageName, tag)

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
