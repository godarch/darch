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

	// Let's get the ID of our current booted image, so we don't delete it.
	currentBootID := ""
	{
		currentBootedImage, err := session.GetCurrentBootedImage()
		if err != nil {
			currentBootID = currentBootedImage.ID
		}
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
			// It isn't in database, but are we currently booting it?
			if liveImage == currentBootID {
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
