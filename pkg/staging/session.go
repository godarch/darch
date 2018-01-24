package staging

import (
	"os"

	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/utils"
)

type Session struct {
	imageStore reference.Store
	imagesDir  string
}

// NewSession creates a new session
func NewSession() (*Session, error) {
	imageStore, err := reference.NewReferenceStore(DefaultStagingImagesFile)
	if err != nil {
		return nil, err
	}

	if !utils.DirectoryExists(DefaultStagingDirectoryImages) {
		err = os.MkdirAll(DefaultStagingDirectoryImages, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &Session{
		imageStore: imageStore,
		imagesDir:  DefaultStagingDirectoryImages,
	}, nil
}
