package main

import (
	"fmt"
	"os"

	"github.com/godarch/darch/pkg/cmd/darch/commands/helpers"
	"github.com/godarch/darch/pkg/cmd/darch/commands/hooks"
	"github.com/godarch/darch/pkg/cmd/darch/commands/images"
	"github.com/godarch/darch/pkg/cmd/darch/commands/recipes"
	"github.com/godarch/darch/pkg/cmd/darch/commands/stage"

	"github.com/urfave/cli"
)

// GitCommit The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version The main version number that is being run at the moment.
var Version = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."
	app.Version = Version
	app.HideVersion = true
	app.Commands = []cli.Command{
		images.Command,
		recipes.Command,
		stage.Command,
		helpers.Command,
		hooks.Command,
		{
			Name:  "version",
			Usage: "Print version information about darch.",
			Action: func(c *cli.Context) error {
				fmt.Printf("version %s\n", Version)
				fmt.Printf("commit %s\n", GitCommit)
				return nil
			},
		},
	}

	markdownDoc := false
	if len(os.Args) >= 2 {
		if os.Args[1] == "markdown" {
			// We are running this command to generate documentation!
			os.Args = append(os.Args[:1], os.Args[2:]...)
			markdownDoc = true
		}
	}

	if markdownDoc {
		for subCommandIndex, subCommand := range app.Commands {
			subCommand.Action = markdownAction
			app.Commands[subCommandIndex] = walkCommandsForDocumentation(subCommand)
		}
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "darch: %s\n", err)
		os.Exit(1)
	}
}

func markdownAction(clicontext *cli.Context) error {
	fmt.Printf("**%s**\n", clicontext.Command.HelpName)
	fmt.Println("")

	if len(clicontext.Command.Description) > 0 {
		fmt.Println("## Description")
		fmt.Println("")
		fmt.Println(clicontext.Command.Description)
		fmt.Println("")
	}

	fmt.Println("## Usage")
	fmt.Println("")
	fmt.Printf(clicontext.Command.HelpName)
	if len(clicontext.Command.ArgsUsage) > 0 {
		fmt.Printf(" %s\n", clicontext.Command.ArgsUsage)
	} else {
		fmt.Println("")
	}
	fmt.Println("")

	return nil
}

func walkCommandsForDocumentation(cmd cli.Command) cli.Command {
	cmd.Action = markdownAction
	for subCommandIndex, subCommand := range cmd.Subcommands {
		cmd.Subcommands[subCommandIndex] = walkCommandsForDocumentation(subCommand)
	}
	return cmd
}
