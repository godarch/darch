package utils

import (
	"fmt"
	"os"
	"path"
)

// ExpandPath Expands the given path to an absolute directory
func ExpandPath(pathToExpand string) string {
	if !path.IsAbs(pathToExpand) {
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Getwd failed: %s", err))
		}
		return path.Clean(path.Join(wd, pathToExpand))
	}
	return pathToExpand
}

// DirectoryExists Returns true if the given path is a directory, and it exists.
func DirectoryExists(directory string) bool {
	if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
		return true
	}
	return false
}

// FileExists Returns true if the given path is a file, and it exists.
func FileExists(directory string) bool {
	if stat, err := os.Stat(directory); err == nil && !stat.IsDir() {
		return true
	}
	return false
}
