package stage

import (
	"github.com/godarch/darch/pkg/cmd/darch/commands/stage/grub"
	"github.com/urfave/cli"
)

var (
	// Command is the cli command for managing content
	Command = cli.Command{
		Name:  "stage",
		Usage: "manage the stage",
		Subcommands: cli.Commands{
			listCommand,
			uploadCommand,
			removeCommand,
			tagCommand,
			runHooksCommand,
			syncBootloaderCommand,
			grub.Command,
		},
	}
)
