package staging

import (
	"fmt"
	"github.com/pauldotknopf/darch/pkg/block"
	"github.com/pauldotknopf/darch/pkg/grub"
	"io"
	"os"
	"path"
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
	}, os.Stdout)
}
