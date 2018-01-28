package commands

import (
	"fmt"
	"os/user"
)

// CheckForRoot Makes sure the running user is root.
func CheckForRoot() error {
	current, err := user.Current()
	if err != nil {
		return err
	}
	if current.Uid != "0" {
		return fmt.Errorf("you must be root")
	}
	return nil
}
