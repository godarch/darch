package gotree

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

/*GTStructure Structure to output print */
type GTStructure struct {
	Name  string
	Items []GTStructure
}

func StringTree(object GTStructure) (result string) {
	result += object.Name + "\n"
	var spaces []bool
	result += stringObjItems(object.Items, spaces)
	return
}

func stringLine(name string, spaces []bool, last bool) (result string) {
	for _, space := range spaces {
		if space {
			result += "    "
		} else {
			result += "│   "
		}
	}

	indicator := "├── "
	if last {
		indicator = "└── "
	}

	result += indicator + name + "\n"
	return
}

func stringObjItems(items []GTStructure, spaces []bool) (result string) {
	for i, f := range items {
		last := (i >= len(items)-1)
		result += stringLine(f.Name, spaces, last)
		if len(f.Items) > 0 {
			spacesChild := append(spaces, last)
			result += stringObjItems(f.Items, spacesChild)
		}
	}
	return
}

/*PrintTree - Print the tree in console */
func PrintTree(object GTStructure) {
	fmt.Println(StringTree(object))
}

/*ReadFolder - Read a folder and return the generated object */
func ReadFolder(directory string) GTStructure {

	var parent GTStructure

	parent.Name = directory
	parent.Items = createGTReadFolder(directory)

	return parent
}

func createGTReadFolder(directory string) []GTStructure {

	var items []GTStructure
	files, _ := ioutil.ReadDir(directory)

	for _, f := range files {

		var child GTStructure
		child.Name = f.Name()

		if f.IsDir() {
			newDirectory := filepath.Join(directory, f.Name())
			child.Items = createGTReadFolder(newDirectory)
		}

		items = append(items, child)
	}
	return items
}
