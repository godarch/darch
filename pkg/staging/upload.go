package staging

import (
	"fmt"
	"os"
	"path"

	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/utils"
)

// UploadDirectoryWithMove Moves (not copy) a directory to staging for boot.
func (session *Session) UploadDirectoryWithMove(imageDir string, imageRef reference.ImageRef, force bool) error {
	// If we aren't forcing this upload, we don't intend to overwrite images already on the stage.
	// So, let's do a quick check to see if it exists.
	if !force {
		_, err := session.imageStore.Get(imageRef)
		if err == nil {
			// We got a valid id for this image, which means it's already staged.
			return fmt.Errorf("image %s already exists in stage", imageRef.FullName())
		}
	}

	img, err := parseImageDir(imageDir)
	if err != nil {
		return err
	}

	// TODO: Need to make sure this doesn't exist. It likely wont.
	newID := utils.NewID()
	newDir := path.Join(session.imagesDir, newID)
	if err != nil {
		return err
	}
	err = os.Rename(img.Dir, newDir)
	if err != nil {
		return err
	}

	img.Dir = newDir

	err = session.imageStore.AddTag(imageRef, newID, force)
	if err != nil {
		// Since we couldn't store this image in database, let's remove the directory.
		os.RemoveAll(newDir)
		return err
	}

	return nil
}
