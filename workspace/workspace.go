package workspace

import (
	"os"
	"path"

	"github.com/pauldotknopf/darch/utils"
)

type Workspace struct {
	Path string
}

func NewWorkspace(tmpDir string) (Workspace, error) {
	var destination = ""
	var found = false

	for !found {
		destination = path.Join(tmpDir, utils.NewID())
		err := os.Mkdir(destination, os.ModePerm)
		if err == nil {
			found = true
		}
	}

	return Workspace{
		Path: destination,
	}, nil
}

func DestroyWorkspace(workspace Workspace) error {
	return os.RemoveAll(workspace.Path)
}
