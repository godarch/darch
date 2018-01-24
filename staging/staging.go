package staging

import (
	"path"
)

var (
	// DefaultStagingDirectory The location where staging happens
	DefaultStagingDirectory = "/var/lib/darch/stage"
	// DefaultStagingDirectoryImages The location to the staged images live.
	DefaultStagingDirectoryImages = path.Join(DefaultStagingDirectory, "live")
	// DefaultStagingDirectoryTmp Temp directory for staging stuff (extracting images, running hooks, etc)
	DefaultStagingDirectoryTmp = path.Join(DefaultStagingDirectory, "tmp")
	// DefaultStagingImagesFile File where our staged images information lives.
	DefaultStagingImagesFile = path.Join(DefaultStagingDirectory, "images.json")
)

// GetAllStaged Get all the staged items in the given directory.
func (session *Session) GetAllStaged() ([]StagedImageNamed, error) {
	result := []StagedImageNamed{}

	associations, err := session.imageStore.AllImages()
	if err != nil {
		return result, err
	}

	for _, association := range associations {
		imageDir := path.Join(DefaultStagingDirectoryImages, association.ID)
		image, err := parseImageDir(imageDir)
		if err != nil {
			return result, err
		}
		result = append(result, StagedImageNamed{
			StagedImage: image,
			Ref:         association.Ref,
		})
	}

	return result, nil
}
