package staging

import (
	"fmt"
	"os"
	"path"

	"github.com/pauldotknopf/darch/reference"
	"github.com/pauldotknopf/darch/utils"
)

// UploadDirectoryWithMove Moves (not copy) a directory to staging for boot.
func UploadDirectoryWithMove(imageDir, imageName string) error {
	imageRef, err := reference.ParseImage(imageName)
	if err != nil {
		return err
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
	newDir := path.Join(DefaultStagingDirectoryImages, utils.NewID())
	if err != nil {
		return err
	}
	err = os.Rename(img.Dir, newDir)
	if err != nil {
		return err
	}

	img.Dir = newDir

	// TODO: Save this image in repository.
	fmt.Println(imageRef.FullName())

	return nil
}
