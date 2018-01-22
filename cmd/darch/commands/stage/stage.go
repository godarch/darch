package stage

import (
	"fmt"
	"os/user"

	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:  "stage",
		Usage: "manage the stage",
		Subcommands: cli.Commands{
			listCommand,
		},
	}
)

func checkForRoot() error {
	current, err := user.Current()
	if err != nil {
		return err
	}
	if current.Uid != "0" {
		return fmt.Errorf("you must be root")
	}
	return nil
}
