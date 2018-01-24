package hooks

import (
	"fmt"

	"github.com/pauldotknopf/darch/hooks"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list hooks",
	Action: func(clicontext *cli.Context) error {

		hooks, err := hooks.GetHooks()
		if err != nil {
			return err
		}

		for _, hook := range hooks {
			fmt.Println(hook.Name)
		}

		return nil
	},
}
