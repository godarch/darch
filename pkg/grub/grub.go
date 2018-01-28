package grub

import (
	"bufio"
	"fmt"
	"github.com/openconfig/goyang/pkg/indent"
	"io"
	"os/exec"
)

// PrepareAccessToDevice Writes out the required info to access a device in grub.
// ---------
// insmod part_gpt
// insmod ext2
// set root='hd3,gpt2'
// if [ x$feature_platform_search_hint = xy ]; then
//   search --no-floppy --fs-uuid --set=root --hint-bios=hd3,gpt2 --hint-efi=hd3,gpt2 --hint-baremetal=ahci3,gpt2  a85bf1c9-59a1-4dba-9df9-6fbbfa03466c
// else
//   search --no-floppy --fs-uuid --set=root a85bf1c9-59a1-4dba-9df9-6fbbfa03466c
// fi
// ---------
func PrepareAccessToDevice(device string, output io.Writer) error {
	if len(device) == 0 {
		return fmt.Errorf("device is required")
	}
	result, err := runCommand("/usr/bin/env",
		"bash",
		"-c",
		fmt.Sprintf(". /usr/share/grub/grub-mkconfig_lib && prepare_grub_to_access_device %s", device))
	if err != nil {
		return err
	}

	for _, line := range result {
		_, err = io.WriteString(output, line+"\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadLinux Prints the loading of the kernel and initramfs.
// ---
// echo	'Loading Linux linux ...'
// linux	/var/lib/darch/stage/live/h42milg9c1r2jqr02ch8q1r3x/vmlinuz-linux root=UUID=a85bf1c9-59a1-4dba-9df9-6fbbfa03466c rw  quiet darch_rootfs=rootfs.squash darch_dir=UUID=a85bf1c9-59a1-4dba-9df9-6fbbfa03466c:/var/lib/darch/stage/live/h42milg9c1r2jqr02ch8q1r3x
// echo	'Loading initial ramdisk ...'
// initrd  /var/lib/darch/stage/live/h42milg9c1r2jqr02ch8q1r3x/initramfs-linux.img
// ---
func LoadLinux(kernelPath string, kernelCommand string, initrdPath string, output io.Writer) error {
	if len(kernelPath) > 0 {
		_, err := io.WriteString(output, "echo 'Loading kernel...'\n")
		if err != nil {
			return err
		}
		_, err = io.WriteString(output, fmt.Sprintf("linux %s", kernelPath))
		if len(kernelCommand) > 0 {
			_, err = io.WriteString(output, fmt.Sprintf(" %s\n", kernelCommand))
			if err != nil {
				return err
			}
		} else {
			_, err = io.WriteString(output, "\n")
			if err != nil {
				return err
			}
		}
	}
	if len(initrdPath) > 0 {
		_, err := io.WriteString(output, "echo 'Loading initial ramdisk...'\n")
		if err != nil {
			return err
		}
		_, err = io.WriteString(output, fmt.Sprintf("initrd %s\n", initrdPath))
		if err != nil {
			return err
		}
	}
	return nil
}

// MenuEntry Generates a menu entry, with a callback to write the contents of the menu entry.
func MenuEntry(name string, contents func(w io.Writer) error, output io.Writer) error {
	if len(name) == 0 {
		return fmt.Errorf("menu entry name required")
	}
	_, err := io.WriteString(output, fmt.Sprintf("menuentry '%s' {\n", name))
	if err != nil {
		return err
	}
	tw := indent.NewWriter(output, "  ")
	err = contents(tw)
	if err != nil {
		return err
	}
	_, err = io.WriteString(output, "}\n")
	if err != nil {
		return err
	}
	return nil
}

func runCommand(name string, args ...string) ([]string, error) {
	cmd := exec.Command(name, args...)
	cmdOut, _ := cmd.StdoutPipe()

	result := []string{}

	err := cmd.Start()
	if err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(cmdOut)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		return result, err
	}

	return result, nil
}
