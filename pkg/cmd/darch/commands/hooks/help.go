package hooks

import (
	"github.com/pauldotknopf/darch/pkg/hooks"
	"github.com/urfave/cli"
)

var helpCommand = cli.Command{
	Name:      "help",
	Usage:     "print help about a hook",
	ArgsUsage: "<hook>",
	Action: func(clicontext *cli.Context) error {
		var (
			hookName = clicontext.Args().First()
		)

		hook, err := hooks.GetHook(hookName)
		if err != nil {
			return err
		}

		return hooks.PrintHookHelp(hook)
	},
}
