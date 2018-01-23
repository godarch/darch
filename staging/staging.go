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
)

// GetAllStaged Get all the staged items in the given directory.
func GetAllStaged() (map[string]StagedImageNamed, error) {
	result := make(map[string]StagedImageNamed, 0)
	return result, nil
}
