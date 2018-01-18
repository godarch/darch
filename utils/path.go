package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
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

// GetChildDirectories Gets child directories in the given directory
func GetChildDirectories(path string) ([]string, error) {
	directories := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			if !strings.HasPrefix(f.Name(), ".") {
				directories = append(directories, f.Name())
			}
		}
	}
	sort.Strings(directories)
	return directories, nil
}
