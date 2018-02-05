package block

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"path"
	"strings"
)

// GetBlockDeviceForPath Get the block device for the given path
func GetBlockDeviceForPath(path string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("path required")
	}

	result, err := runCommand("df", "--output=source", path)
	if err != nil {
		return "", err
	}

	if len(result) == 2 {
		if result[1] == "-" {
			return "", fmt.Errorf("error getting block device for %s, possibly not on block device", path)
		}
		return result[1], nil
	}

	return "", errors.Wrap(fmt.Errorf(strings.Join(result, " ")), "error running df")
}

// GetUUIDForBlockDevice Get the UUID for the given block device.
func GetUUIDForBlockDevice(blockDevice string) (string, error) {
	if len(blockDevice) == 0 {
		return "", fmt.Errorf("block device required")
	}

	result, err := runCommand("blkid", blockDevice)
	if err != nil {
		return "", err
	}

	if len(result) == 1 {
		values := strings.Split(result[0], " ")
		for _, value := range values {
			if strings.HasPrefix(value, "UUID=\"") {
				return value[6 : len(value)-1], nil
			}
		}
	}

	return "", errors.Wrap(fmt.Errorf(strings.Join(result, " ")), "error running blkid")
}

// GetPathRelativeToBlockDevice Give a full path to your system, and it will return it's path, relative to the device it is hosted on.
// For example, if a device is mounted on "/test/mount" and you invoke this method with "/test/mount/with/this/file",
// then you will get /with/this/file returned.
func GetPathRelativeToBlockDevice(p string) (string, error) {
	if len(p) == 0 {
		return "", fmt.Errorf("path required")
	}

	result, err := runCommand("df", p, "--output=target")
	if err != nil {
		return "", err
	}

	if len(result) == 2 && result[0] == "Mounted on" {
		return path.Join("/", p[len(result[1]):]), nil
	}

	return "", nil
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
