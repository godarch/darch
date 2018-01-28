package grub

import (
	"io"
	"os"
	"testing"
)

func TestAccess(t *testing.T) {
	t.Skip()
	err := PrepareAccessToDevice("/dev/sdd2", os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKernel(t *testing.T) {
	t.Skip()
	err := LoadLinux("/path/to/vmlinuz", "kernel-command-line=test", "/path/to/initrd.img", os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMenuEntry(t *testing.T) {
	t.Skip()
	err := MenuEntry("test entry", func(w io.Writer) error {
		err := PrepareAccessToDevice("/dev/sdd2", w)
		if err != nil {
			return err
		}
		err = LoadLinux("/path/to/vmlinuz", "kernel-command-line=test", "/path/to/initrd.img", w)
		if err != nil {
			return err
		}
		return nil
	}, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
