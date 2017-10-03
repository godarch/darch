package images

import (
	"fmt"
	"path"

	"../utils"
)

// Image A struct representing an image to be built.
type Image struct {
	Name     string
	ImageDir string
	Inherits []string
}

type imageConfiguration struct {
	name     string
	inherits string
}

// BuildDefinition Parse an image from the file system
func BuildDefinition(imageName string, imageDir string) (Image, error) {

	image := Image{}

	image.ImageDir = utils.ExpandPath(imageDir)
	image.Name = imageName

	if !utils.DirectoryExists(imageDir) {
		return Image{}, fmt.Errorf("Image directory %s doesn't exist", imageDir)
	}

	if len(imageName) == 0 {
		return image, fmt.Errorf("An image must be provided")
	}

	if !utils.DirectoryExists(imageDir + "/" + imageName) {
		return image, fmt.Errorf("The image %s doesn't exist in %s", imageName, imageDir)
	}

	return image, nil
}

func loadImageConfiguration(image Image) (imageConfiguration, error) {
	imageConfigurationPath := path.Join(image.ImageDir, image.Name, "config.json")
	imageConfiguration := imageConfiguration{}

	if !utils.FileExists(imageConfigurationPath) {
		return imageConfiguration, fmt.Errorf("No configuration file exists at %s", imageConfigurationPath)
	}

	//json, err := ioutil.ReadFile(imageConfigurationPath)

	//if err != nil {
	//	return imageConfiguration, err
	//}

	//dec := json.NewDecoder(strings(json))

	// for {
	// 	var m Message
	// 	if err := dec.Decode(&m); err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%s: %s\n", m.Name, m.Text)
	// }

	return imageConfiguration, nil
}
