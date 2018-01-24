package utils

import "os"

// FileExists Returns true if the given path is a file, and it exists.
func FileExists(file string) bool {
	if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
		return true
	}
	return false
}
