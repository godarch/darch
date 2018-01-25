package helpers

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pauldotknopf/darch/pkg/utils"
	"github.com/urfave/cli"
)

var globFirstMatchCommand = cli.Command{
	Name:  "glob-config-first-match",
	Usage: "print the first matching value from a glob config",
	Action: func(clicontext *cli.Context) error {
		var (
			file  = clicontext.Args().Get(0)
			value = clicontext.Args().Get(1)
		)

		if len(file) == 0 {
			return fmt.Errorf("no file given")
		}

		if len(value) == 0 {
			return fmt.Errorf("no value given")
		}

		lines, err := utils.GetFileLines(file)
		if err != nil {
			return err
		}

		for _, line := range lines {
			delimiterPosition := strings.LastIndex(line, "=")
			if delimiterPosition == -1 {
				return fmt.Errorf("invalid entry \"%s\"", line)
			}

			globPattern := line[:delimiterPosition]
			globValue := line[delimiterPosition+1:]

			if len(globPattern) == 0 {
				return fmt.Errorf("invalid entry \"%s\"", line)
			}

			g := glob.MustCompile(globPattern)

			if g.Match(value) {
				fmt.Println(globValue)
				return nil
			}
		}

		return fmt.Errorf("No matching entries")
	},
}
