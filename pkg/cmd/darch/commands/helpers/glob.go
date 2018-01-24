package helpers

import (
	"fmt"

	"github.com/gobwas/glob"
	"github.com/urfave/cli"
)

var globCommand = cli.Command{
	Name:  "glob",
	Usage: "test pattern/input using globbing",
	Action: func(clicontext *cli.Context) error {
		var (
			pattern = clicontext.Args().Get(0)
			value   = clicontext.Args().Get(1)
		)

		if len(pattern) == 0 {
			return fmt.Errorf("no pattern given")
		}

		if len(value) == 0 {
			return fmt.Errorf("no value given")
		}

		g := glob.MustCompile(pattern)
		if !g.Match(value) {
			return fmt.Errorf("not a match")
		}

		fmt.Println("Match!")

		return nil
	},
}
