package images

import (
	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:  "images",
		Usage: "manage images",
		Subcommands: cli.Commands{
			pullCommand,
			listCommand,
			tagCommand,
			removeCommand,
		},
	}
)
