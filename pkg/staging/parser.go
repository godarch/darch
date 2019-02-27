package staging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/godarch/darch/pkg/reference"
	"github.com/godarch/darch/pkg/utils"
)

// StagedImage A staged image
type StagedImage struct {
	Dir           string
	Kernel        string
	KernelParams  string
	InitRAMFS     string
	RootFS        string
	NoDoubleMount bool
	CreationTime  time.Time
}

// StagedImageNamed A StagedImage with a name and tag
type StagedImageNamed struct {
	StagedImage
	Ref reference.ImageRef
	ID  string
}

type stagedImageConfiguration struct {
	Kernel        string `json:"kernel"`
	KernelParams  string `json:"kernelparams"`
	InitRAMFS     string `json:"initramfs"`
	RootFS        string `json:"rootfs"`
	NoDoubleMount bool   `json:"nodoublemount"`
}

// ParseImageDir Parses an image directory, and also validates it.
func parseImageDir(imageDir string) (StagedImage, error) {
	result := StagedImage{}
	result.Dir = imageDir

	if !utils.DirectoryExists(imageDir) {
		return result, fmt.Errorf("invalid image directory")
	}

	config, err := loadStagedImageConfiguration(path.Join(imageDir, "image.json"))
	if err != nil {
		return result, err
	}

	if len(config.InitRAMFS) == 0 {
		return result, fmt.Errorf("initramfs was empty")
	}

	if len(config.Kernel) == 0 {
		return result, fmt.Errorf("kernel was empty")
	}

	if len(config.RootFS) == 0 {
		return result, fmt.Errorf("rootfs was empty")
	}

	if !utils.FileExists(path.Join(imageDir, config.InitRAMFS)) {
		return result, fmt.Errorf("initramfs was invalid")
	}

	if !utils.FileExists(path.Join(imageDir, config.Kernel)) {
		return result, fmt.Errorf("kernel was invalid")
	}

	if !utils.FileExists(path.Join(imageDir, config.RootFS)) {
		return result, fmt.Errorf("rootfs was invalid")
	}

	stat, err := os.Stat(imageDir)
	if err != nil {
		return result, err
	}

	result.InitRAMFS = config.InitRAMFS
	result.Kernel = config.Kernel
	result.KernelParams = config.KernelParams
	result.RootFS = config.RootFS
	result.NoDoubleMount = config.NoDoubleMount
	result.CreationTime = stat.ModTime()

	return result, nil
}

func loadStagedImageConfiguration(file string) (stagedImageConfiguration, error) {

	result := stagedImageConfiguration{}

	if !utils.FileExists(file) {
		return result, fmt.Errorf("no configuration file exists at %s", file)
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
