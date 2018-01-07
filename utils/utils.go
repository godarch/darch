package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

// RemoveDuplicates Remove duplicate items from an array.
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// Reverse Reverses the array
func Reverse(elements []string) []string {
	for i := len(elements)/2 - 1; i >= 0; i-- {
		opp := len(elements) - 1 - i
		elements[i], elements[opp] = elements[opp], elements[i]
	}
	return elements
}

// Contains Returns true of array contains the given string
func Contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

// DirectoryExists Returns true if the given path is a directory, and it exists.
func DirectoryExists(directory string) bool {
	if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
		return true
	}
	return false
}

// FileExists Returns true if the given path is a file, and it exists.
func FileExists(file string) bool {
	if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
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

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// ConvertVariableStringsToMap Converts an array of "KEY=VALUE" strings
// to a hash table of "hash["KEY"] = "VALUE""".
func ConvertVariableStringsToMap(input []string) (map[string]string, error) {
	m := make(map[string]string)

	for _, value := range input {
		if len(value) == 0 {
			continue
		}
		split := strings.Split(value, "=")
		if len(split) > 2 {
			return nil, fmt.Errorf("Invalid environment value %s", value)
		}
		if len(split) == 2 {
			m[split[0]] = split[1]
			continue
		}
		if len(split) == 1 {
			m[split[0]] = ""
		}
	}

	return m, nil
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
