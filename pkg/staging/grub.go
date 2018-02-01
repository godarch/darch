package staging

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/godarch/darch/pkg/block"
	"github.com/godarch/darch/pkg/grub"
	"io"
	"os"
	"path"
)

var (
	// DefaultGrubConfigPath The location where the grub.cfg for darch is stored.
	DefaultGrubConfigPath = "/etc/darch/grub.cfg"
)

// PrintGrubMenuEntry Print the grub entry for the given staged image.
func (session *Session) PrintGrubMenuEntry(stagedImage StagedImageNamed, output io.Writer) error {
	device, err := block.GetBlockDeviceForPath(stagedImage.Dir)
	if err != nil {
		return err
	}
	relPathTodevice, err := block.GetPathRelativeToBlockDevice(stagedImage.Dir)
	if err != nil {
		return err
	}
	uuid, err := block.GetUUIDForBlockDevice(device)
	if err != nil {
		return err
	}
	commandLine := fmt.Sprintf("darch_rootfs=%s darch_dir=UUID=%s:%s", stagedImage.RootFS, uuid, relPathTodevice)

	return grub.MenuEntry(fmt.Sprintf("Darch - %s", stagedImage.Ref.FullName()), func(w io.Writer) error {
		err := grub.PrepareAccessToDevice(device, w)
		if err != nil {
			return err
		}
		err = grub.LoadLinux(path.Join(relPathTodevice, stagedImage.Kernel),
			commandLine,
			path.Join(relPathTodevice, stagedImage.InitRAMFS),
			w)
		if err != nil {
			return err
		}
		return nil
	}, output)
}

// SyncBootloader Updates the /etc/darch/grub.cfg to represent the current stage.
func (session *Session) SyncBootloader() error {
	allImages, err := session.GetAllStaged()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	for _, image := range allImages {
		err = session.PrintGrubMenuEntry(image, w)
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return ioutils.AtomicWriteFile(DefaultGrubConfigPath, b.Bytes(), os.ModePerm)
}
