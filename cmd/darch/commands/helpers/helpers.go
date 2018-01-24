package helpers

import (
	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:   "helpers",
		Usage:  "helpers tasks for various functions/hooks",
		Hidden: true,
		Subcommands: cli.Commands{
			globCommand,
		},
	}
)
