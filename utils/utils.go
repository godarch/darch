package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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

// CleanDirectory Wipes all the data within a folder
func CleanDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
