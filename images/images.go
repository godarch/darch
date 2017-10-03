package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"../utils"
)

// Image A struct representing an image to be built.
type ImageDefinition struct {
	Name     string
	ImageDir string
	Inherits []string
}

type imageConfiguration struct {
	Inherits string `json:"inherits"`
}

// BuildDefinition Parse an image from the file system
func BuildDefinition(imageName string, imageDir string) (*ImageDefinition, error) {

	image := ImageDefinition{}

	image.ImageDir = utils.ExpandPath(imageDir)
	image.Name = imageName

	if !utils.DirectoryExists(imageDir) {
		return nil, fmt.Errorf("Image directory %s doesn't exist", imageDir)
	}

	if len(imageName) == 0 {
		return nil, fmt.Errorf("An image must be provided")
	}

	if !utils.DirectoryExists(imageDir + "/" + imageName) {
		return nil, fmt.Errorf("The image %s doesn't exist in %s", imageName, imageDir)
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

func loadImageConfiguration(image ImageDefinition) (*imageConfiguration, error) {
	imageConfigurationPath := path.Join(image.ImageDir, image.Name, "config.json")
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
