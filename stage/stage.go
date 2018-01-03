package stage

import (
	"fmt"
	"time"
)

type StagedItemTag struct {
	Name          string
	StagedOn      time.Time
	Path          string
	BootKernel    string
	BootInitRAMFS string
	BootRootFs    string
}

// StagedItem A struct representing a a staged item
type StagedItem struct {
	Name          string
	Tag           string
	StagedOn      time.Time
	Path          string
	BootKernel    string
	BootInitRAMFS string
	BootRootFs    string
	Tags          []StagedItemTag
}

// GetAllStaged Get all the staged items in the given directory.
func GetAllStaged(stagedDirectory string) (map[string]StagedItem, error) {
	if len(stagedDirectory) == 0 {
		return nil, fmt.Errorf("A staging directory must be provided")
	}
}
