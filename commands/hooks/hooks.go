package hooks

import (
	"fmt"

	"strconv"

	"../../hooks"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

func listCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List hooks currently installed.",
		Action: func(c *cli.Context) error {
			err := list()
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
		Name:  "hooks",
		Usage: "Commands manage/view hooks.",
		Subcommands: []cli.Command{
			listCommand(),
		},
	}
}

func list() error {

	hooks, err := hooks.GetHooks()

	if err != nil {
		return err
	}

	result := []string{
		"Name | Path | ExecutionOrder",
	}

	for _, hook := range hooks {
		result = append(result, hook.Name+" | "+hook.Path+" | "+strconv.Itoa(hook.ExecutionOrder))
	}

	fmt.Println(columnize.SimpleFormat(result))

	return nil
}
