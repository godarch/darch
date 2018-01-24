package hooks

import (
	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:  "hooks",
		Usage: "view hooks on the system",
		Subcommands: cli.Commands{
			listCommand,
			helpCommand,
			detailsCommand,
		},
	}
)
