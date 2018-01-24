package staging

import (
	"fmt"
	"os"
	"path"

	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/utils"
)

// UploadDirectoryWithMove Moves (not copy) a directory to staging for boot.
func UploadDirectoryWithMove(imageDir string, imageRef reference.ImageRef, force bool) error {
	imageStore, err := reference.NewReferenceStore(DefaultStagingImagesFile)
	if err != nil {
		return err
	}

	// If we aren't forcing this upload, we don't intend to overwrite images already on the stage.
	// So, let's do a quick check to see if it exists.
	if !force {
		_, err := imageStore.Get(imageRef)
		if err == nil {
			// We got a valid id for this image, which means it's already staged.
			return fmt.Errorf("image %s already exists in stage", imageRef.FullName())
		}
	}

	img, err := ParseImageDir(imageDir)
	if err != nil {
		return err
	}

	if !utils.DirectoryExists(DefaultStagingDirectoryImages) {
		err = os.MkdirAll(DefaultStagingDirectoryImages, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// TODO: Need to make sure this doesn't exist. It likely wont.
	newID := utils.NewID()
	newDir := path.Join(DefaultStagingDirectoryImages, newID)
	if err != nil {
		return err
	}
	err = os.Rename(img.Dir, newDir)
	if err != nil {
		return err
	}

	img.Dir = newDir

	err = imageStore.AddTag(imageRef, newID, force)
	if err != nil {
		// Since we couldn't store this image in database, let's remove the directory.
		os.RemoveAll(newDir)
		return err
	}

	return nil
}
