package workspace

import (
	"io/ioutil"
	"os"

	"github.com/pauldotknopf/darch/pkg/utils"
)

// Workspace The workspace.
type Workspace struct {
	Path      string
	destroyed bool
}

// NewWorkspace Create a new temporary workspace.
func NewWorkspace(tmpDir string) (Workspace, error) {
	if len(tmpDir) > 0 {
		if !utils.DirectoryExists(tmpDir) {
			err := os.MkdirAll(tmpDir, os.ModePerm)
			if err != nil {
				return Workspace{}, err
			}
		}
	}

	path, err := ioutil.TempDir(tmpDir, utils.NewID())

	if err != nil {
		return Workspace{}, err
	}

	return Workspace{
		Path:      path,
		destroyed: false,
	}, nil
}

// Destroy Destroy the workspace.
func (workspace *Workspace) Destroy() error {
	if workspace.destroyed {
		return nil
	}
	err := os.RemoveAll(workspace.Path)
	if err != nil {
		return err
	}
	workspace.destroyed = true
	return nil
}

// MarkDestroyed If called, the Destroy method will do nothing.
func (workspace *Workspace) MarkDestroyed() {
	workspace.destroyed = true
}
