package staging

import (
	"os"
	"path"

	"github.com/godarch/darch/pkg/utils"
)

// Clean goes through all the images in the live directory and deletes them
// if there isn't a references in images.json.
func (session *Session) Clean() error {
	liveImages, err := utils.GetChildDirectories(DefaultStagingDirectoryImages)
	if err != nil {
		return err
	}

	databaseImages, err := session.imageStore.AllImages()
	if err != nil {
		return err
	}

	for _, liveImage := range liveImages {
		found := false
		for _, databaseImage := range databaseImages {
			if databaseImage.ID == liveImage {
				found = true
			}
		}
		if !found {
			err = os.RemoveAll(path.Join(DefaultStagingDirectoryImages, liveImage))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
