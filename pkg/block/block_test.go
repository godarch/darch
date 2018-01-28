package block

import (
	"testing"
)

func TestBlockDevice(t *testing.T) {
	t.Skip()
	device, err := GetBlockDeviceForPath("/usr/bin/darch")
	if err != nil {
		t.Fatal(err)
	}
	if len(device) == 0 {
		t.Fatal("empty device")
	}
}

func TestUUID(t *testing.T) {
	t.Skip()
	uuid, err := GetUUIDForBlockDevice("/dev/sdd2")
	if err != nil {
		t.Fatal(err)
	}
	if len(uuid) == 0 {
		t.Fatal("empty uuid")
	}
}

func TestRelative(t *testing.T) {
	t.Skip()
	rel, err := GetPathRelativeToBlockDevice("/boot/grub/grub.cfg")
	if err != nil {
		t.Fatal(err)
	}
	if len(rel) == 0 {
		t.Fatal("empty rel")
	}
}
