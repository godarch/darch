package helpers

import (
	"fmt"
	"strings"

	"../../utils"
	"github.com/gobwas/glob"
	"github.com/urfave/cli"
)

func globCommand() cli.Command {
	return cli.Command{
		Name:      "glob",
		Usage:     "Test a pattern/input using globbing.",
		ArgsUsage: "PATTERN VALUE",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 2 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := globPattern(c.Args().First(), c.Args().Get(1))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func globConfigFirstMatchCommand() cli.Command {
	return cli.Command{
		Name:      "glob-config-first-match",
		Usage:     "Print the first matching value from a grub config.",
		ArgsUsage: "CONFIG_FILE VALUE",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 2 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := globPatternConfigFirstMatch(c.Args().First(), c.Args().Get(1))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "helpers",
		Usage: "Random helpers, mostly for hooks, to keep things simple.",
		Subcommands: []cli.Command{
			globCommand(),
			globConfigFirstMatchCommand(),
		},
		Hidden: true,
	}
}

func globPattern(pattern string, value string) error {
	g := glob.MustCompile(pattern)
	if !g.Match(value) {
		return fmt.Errorf("Not a match")
	}
	fmt.Println("Match!")
	return nil
}

func globPatternConfigFirstMatch(file string, value string) error {
	lines, err := utils.GetFileLines(file)
	if err != nil {
		return err
	}

	for _, line := range lines {
		delimiterPosition := strings.LastIndex(line, "=")
		if delimiterPosition == -1 {
			return fmt.Errorf("Invalid entry \"%s\"", line)
		}

		globPattern := line[:delimiterPosition]
		globValue := line[delimiterPosition+1 : len(line)]

		if len(globPattern) == 0 {
			return fmt.Errorf("Invalid entry \"%s\"", line)
		}

		g := glob.MustCompile(globPattern)

		if g.Match(value) {
			fmt.Println(globValue)
			return nil
		}
	}

	return fmt.Errorf("No matching entries")
}
