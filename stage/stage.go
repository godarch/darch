package stage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"../utils"
)

// StagedItemTag A struct representing a tag that is stored.
type StagedItemTag struct {
	Name          string
	StagedOn      time.Time
	Path          string
	BootKernel    string
	BootInitRAMFS string
	BootRootFS    string
}

// StagedItem A struct representing a a staged item.
type StagedItem struct {
	Name string
	Path string
	Tags []StagedItemTag
}

type stagedItemTagConfiguration struct {
	Kernel    string `json:"kernel"`
	InitRAMFS string `json:"initramfs"`
	RootFS    string `json:"rootfs"`
}

// GetAllStaged Get all the staged items in the given directory.
func GetAllStaged(stagedDirectory string) (map[string]StagedItem, error) {

	result := make(map[string]StagedItem, 0)

	if len(stagedDirectory) == 0 {
		return nil, fmt.Errorf("A staging directory must be provided")
	}

	stagedDirectory = utils.ExpandPath(stagedDirectory)

	err := filepath.Walk(stagedDirectory, func(path string, f os.FileInfo, err error) error {
		if f.Name() == "image.json" && !f.IsDir() {
			// This directory is an image.

			stagedItemTagConfiguration, err := loadStagedImageTagConfiguration(path)

			if err != nil {
				return err
			}

			// path = /var/darch/staging/arch/local/image.json
			// directory = /var/darch/staging/arch/local
			// directoryWithoutStaging = arch/local
			directory := filepath.Dir(path)
			directoryWithoutStaging := directory[len(stagedDirectory)+1 : len(directory)]

			imageName := filepath.Dir(directoryWithoutStaging)
			tagName := filepath.Base(directoryWithoutStaging)

			stagedItem := StagedItem{}

			if val, ok := result[imageName]; ok {
				stagedItem = val
			} else {
				stagedItem.Name = imageName
				stagedItem.Path = filepath.Dir(directory)
			}

			stagedItemTag := StagedItemTag{}
			stagedItemTag.Name = tagName
			stagedItemTag.Path = directory
			stagedItemTag.BootKernel = stagedItemTagConfiguration.Kernel
			stagedItemTag.BootInitRAMFS = stagedItemTagConfiguration.InitRAMFS
			stagedItemTag.BootRootFS = stagedItemTagConfiguration.RootFS

			stagedItem.Tags = append(stagedItem.Tags, stagedItemTag)

			result[stagedItem.Name] = stagedItem
		}
		return nil
	})

	return result, err
}

func loadStagedImageTagConfiguration(file string) (stagedItemTagConfiguration, error) {

	result := stagedItemTagConfiguration{}

	if !utils.FileExists(file) {
		return result, fmt.Errorf("No configuration file exists at %s", file)
	}

	jsonData, err := ioutil.ReadFile(file)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(jsonData, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}
