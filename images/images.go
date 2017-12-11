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
	Name             string
	ImageDir         string
	ImagesDir        string
	Inherits         string
	InheritsExternal bool
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

	if strings.HasPrefix(imageConfiguration.Inherits, "external:") {
		image.InheritsExternal = true
		image.Inherits = imageConfiguration.Inherits[len("external:"):len(imageConfiguration.Inherits)]
	} else {
		image.InheritsExternal = false
		image.Inherits = imageConfiguration.Inherits
	}

	return &image, nil
}

func verifyDependencies(imageDefinition ImageDefinition, imageDefinitions map[string]ImageDefinition, currentStack map[string]bool) error {
	if currentStack == nil {
		currentStack = make(map[string]bool, 0)
	}

	if imageDefinition.InheritsExternal {
		// we reached the end, all good!
		return nil
	}

	if _, ok := currentStack[imageDefinition.Inherits]; ok {
		// Cyclical dependency detected!
		return fmt.Errorf("Image %s has a cyclical dependency", imageDefinition.Name)
	}

	// Make this image as having been traversed.
	currentStack[imageDefinition.Name] = true

	if parent, ok := imageDefinitions[imageDefinition.Inherits]; ok {
		return verifyDependencies(parent, imageDefinitions, currentStack)
	}

	return fmt.Errorf("Image defintion %s inherits from %s, which doesn't exist", imageDefinition.Name, imageDefinition.Inherits)
}

// BuildAllDefinitions Return all the images in the image directory
func BuildAllDefinitions(imagesDir string) (map[string]ImageDefinition, error) {
	if len(imagesDir) == 0 {
		return nil, fmt.Errorf("An image directory must be provided")
	}

	imageNames, err := utils.GetChildDirectories(imagesDir)

	if err != nil {
		return nil, err
	}

	imageDefinitions := make(map[string]ImageDefinition, 0)

	for _, imageName := range imageNames {
		imageDefinition, err := BuildDefinition(imageName, imagesDir)
		if err != nil {
			return nil, err
		}
		imageDefinitions[imageName] = *imageDefinition
	}

	// verify dependencies are satisfied and no circular dependencies
	for _, imageDefinition := range imageDefinitions {
		err := verifyDependencies(imageDefinition, imageDefinitions, nil)
		if err != nil {
			return nil, err
		}
	}

	return imageDefinitions, nil
}

func joinStringArrays(array ...[]string) []string {
	result := make([]string, 0)
	for _, item := range array {
		result = append(result, item...)
	}
	return result
}

// BuildImageLayer Run installation scripts on top of another image.
func BuildImageLayer(imageDefinition *ImageDefinition, tags []string, buildPrefix string, packageCache string, environmentVariables map[string]string) error {

	if len(packageCache) > 0 {
		if !utils.DirectoryExists(packageCache) {
			err := os.MkdirAll(packageCache, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	inherits := imageDefinition.Inherits
	if imageDefinition.InheritsExternal {
		// No build prefix for externally referenced image.
	} else {
		inherits = buildPrefix + inherits
	}

	log.Println("Building image " + buildPrefix + imageDefinition.Name + ".")
	log.Println("Using parent image " + inherits + ".")

	// Build the set of arguments that contain the local volumes we are going to mount
	volumeArguments := []string{
		"-v",
		imageDefinition.ImagesDir + ":/images",
	}
	if len(packageCache) > 0 {
		volumeArguments = append(volumeArguments, []string{
			"-v",
			packageCache + ":/packages",
		}...)
	}

	// Build the set of environment variables that we are going to use
	environmentArguements := make([]string, 0)
	for environmentVariableName, environmentVariableValue := range environmentVariables {
		environmentArguements = append(environmentArguements, []string{
			"-e",
			environmentVariableName + "=" + environmentVariableValue,
		}...)
	}

	tmpImageName := "darch-building-" + imageDefinition.Name

	err := runCommand("docker", joinStringArrays(
		[]string{
			"run",
			"-d",
		},
		volumeArguments,
		[]string{
			"--privileged",
			"--name",
			tmpImageName,
			inherits,
		},
	)...)
	if err != nil {
		return err
	}
	// Prep the container
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "/darch-prepare")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	// Run the image scripts
	err = runCommand("docker", joinStringArrays(
		[]string{
			"exec",
			"--privileged",
		},
		environmentArguements,
		[]string{
			tmpImageName,
			"/darch-runimage",
			imageDefinition.Name,
		},
	)...)
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	// Tear the container down
	err = runCommand("docker", "exec", "--privileged", tmpImageName, "/darch-teardown")
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	// Commit the container
	err = runCommand("docker", "commit", tmpImageName, buildPrefix+imageDefinition.Name)
	if err != nil {
		destroyContainer(tmpImageName)
		return err
	}

	// And tag it
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
