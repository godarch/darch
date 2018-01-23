package workspace

import (
	"os"
	"path"

	"github.com/pauldotknopf/darch/utils"
)

// Workspace The workspace.
type Workspace struct {
	Path string
}

// NewWorkspace Create a new temporary workspace.
func NewWorkspace(tmpDir string) (Workspace, error) {
	var destination = ""
	var found = false

	for !found {
		destination = path.Join(tmpDir, utils.NewID())
		err := os.MkdirAll(destination, os.ModePerm)
		if err == nil {
			found = true
		}
	}

	return Workspace{
		Path: destination,
	}, nil
}

// Destroy Destroy the workspace.
func (workspace *Workspace) Destroy() error {
	return os.RemoveAll(workspace.Path)
}
