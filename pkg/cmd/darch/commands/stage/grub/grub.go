package grub

import (
	"github.com/urfave/cli"
)

var (
	// Command The grub helper commands.
	Command = cli.Command{
		Name:   "grub",
		Usage:  "grub helpers for the stage",
		Hidden: true,
		Subcommands: cli.Commands{
			grubMenuEntryCommand,
			grubConfigEntryCommand,
			grubMenuEntriesCommand,
		},
	}
)
