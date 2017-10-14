package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"../utils"
)

// ImageDefinition A struct representing an image to be built.
type ImageDefinition struct {
	Name      string
	ImageDir  string
	ImagesDir string
	Inherits  []string
}

type imageConfiguration struct {
	Inherits string `json:"inherits"`
}

// BuildDefinition Parse an image from the file system
func BuildDefinition(imageName string, imagesDir string) (*ImageDefinition, error) {

	if len(imageName) == 0 {
		return nil, fmt.Errorf("An image must be provided")
	}

	if len(imagesDir) == 0 {
		return nil, fmt.Errorf("An image directory must be provided")
	}

	image := ImageDefinition{}

	image.ImagesDir = utils.ExpandPath(imagesDir)
	image.ImageDir = path.Join(image.ImagesDir, imageName)
	image.Name = imageName

	if !utils.DirectoryExists(image.ImageDir) {
		return nil, fmt.Errorf("Image directory %s doesn't exist", image.ImageDir)
	}

	imageConfiguration, err := loadImageConfiguration(image)

	if err != nil {
		return nil, err
	}

	image.Inherits = []string{
		imageConfiguration.Inherits,
	}

	return &image, nil
}

// BuildImageLayer Run installation scripts on top of another image.
func BuildImageLayer(imageDefinition *ImageDefinition, tags []string, buildPrefix string) error {
	inherits := imageDefinition.Inherits[0]
	if strings.HasPrefix(inherits, "external:") {
		inherits = inherits[len("external:"):len(inherits)]
	} else {
		inherits = buildPrefix + inherits
	}

	log.Println("Building image " + buildPrefix + imageDefinition.Name + ".")
	log.Println("Using parent image " + inherits + ".")

	tmpImageName := "darch-building-" + imageDefinition.Name

	err := runCommand("docker", "run", "-d", "-v", imageDefinition.ImagesDir+":/images", "--privileged", "--name", tmpImageName, inherits)
	if err != nil {
		return err
	}
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "cp", "-rp", "/images", "/root.x86_64/")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "arch-chroot", "/root.x86_64", "/bin/bash", "-c", "cd /images/"+imageDefinition.Name+" && ./script")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "rm", "-r", "/root.x86_64/images")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "rm", "-r", "-f", "/root.x86_64/var/cache/pacman/pkg/")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	err = runCommand("docker", "commit", tmpImageName, buildPrefix+imageDefinition.Name)
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	for _, tag := range tags {
		err = runCommand("docker", "tag", imageDefinition.Name, buildPrefix+imageDefinition.Name+":"+tag)
		if err != nil {
			destroyContainer(tmpImageName)
			return err
		}
	}

	return destroyContainer(tmpImageName)
}

// ExtractImage Extracts an image (with tag) to a specified directory
func ExtractImage(name string, tag string, destination string) error {
	tmpImageName := "darch-extracting-" + strings.Replace(name, "/", "", -1)

	imageName := name
	if len(tag) > 0 {
		imageName = imageName + ":" + tag
	}

	if !utils.DirectoryExists(destination) {
		err := os.MkdirAll(destination, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := utils.CleanDirectory(destination)
	if err != nil {
		return err
	}

	err = runCommand("docker", "run", "-d", "--privileged", "--name", tmpImageName, imageName)
	if err != nil {
		return err
	}

	err = runCommand("docker", "exec", tmpImageName, "mksquashfs", "root.x86_64", "/rootfs.squash")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	err = runCommand("docker", "cp", tmpImageName+":/rootfs.squash", path.Join(destination, "rootfs.squash"))
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	err = runCommand("docker", "cp", tmpImageName+":/root.x86_64/boot/vmlinuz-linux", path.Join(destination, "vmlinuz-linux"))
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	err = runCommand("docker", "cp", tmpImageName+":/root.x86_64/boot/initramfs-linux.img", path.Join(destination, "initramfs-linux.img"))
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	return destroyContainer(tmpImageName)
}

func loadImageConfiguration(image ImageDefinition) (*imageConfiguration, error) {
	imageConfigurationPath := path.Join(image.ImageDir, "config.json")
	imageConfiguration := imageConfiguration{}

	if !utils.FileExists(imageConfigurationPath) {
		return nil, fmt.Errorf("No configuration file exists at %s", imageConfigurationPath)
	}

	jsonData, err := ioutil.ReadFile(imageConfigurationPath)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &imageConfiguration)

	if err != nil {
		return nil, err
	}

	if len(imageConfiguration.Inherits) == 0 {
		return nil, fmt.Errorf("No inherit property given for image %s", image.Name)
	}

	return &imageConfiguration, nil
}

func runCommand(program string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func destroyContainer(containerName string) error {
	err := runCommand("docker", "stop", containerName)
	if err != nil {
		return err
	}

	err = runCommand("docker", "rm", containerName)
	if err != nil {
		return err
	}

	return nil
}
