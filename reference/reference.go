package reference

import (
	"fmt"
	"strings"

	containerdref "github.com/containerd/containerd/reference"
)

type Image struct {
	Name string
	Tag  string
}

// ParseImage Parses a string for image:tag.
func ParseImage(val string) (Image, error) {
	if len(val) == 0 {
		return Image{}, fmt.Errorf("no image name provided")
	}
	result := Image{}
	spec, err := containerdref.Parse(val)
	if err != nil {
		return result, err
	}
	if len(spec.Object) != 0 {
		result.Name = spec.Locator
		result.Tag = spec.Object
	}
	split := strings.Split(val, ":")
	if len(split) > 2 {
		return result, fmt.Errorf("invalid format")
	}
	result.Name = split[0]
	if len(split) == 2 {
		result.Tag = split[1]
	}

	if len(result.Tag) == 0 {
		result.Tag = "latest"
	}
	if len(result.Name) == 0 {
		return result, fmt.Errorf("invalid format")
	}
	return result, nil
}

// FullName Returns image:tag for the image reference.
func (image Image) FullName() string {
	return image.Name + ":" + image.Tag
}
