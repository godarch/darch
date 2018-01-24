package recipes

import (
	"fmt"
	"log"

	"github.com/pauldotknopf/darch/pkg/recipes"
	"github.com/pauldotknopf/darch/pkg/utils"
	"github.com/urfave/cli"
)

var parentsCommand = cli.Command{
	Name:  "parents",
	Usage: "list all the parents of a recipe",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "exclude-external",
		},
		cli.BoolFlag{
			Name: "reverse",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			recipeName      = clicontext.Args().First()
			excludeExternal = clicontext.Bool("exclude-external")
			reverse         = clicontext.Bool("reverse")
		)

		if len(recipeName) == 0 {
			return fmt.Errorf("You must provide a recipe name")
		}

		rs, err := recipes.GetAllRecipes(getRecipesDir(clicontext))
		if err != nil {
			return err
		}

		current, ok := rs[recipeName]
		if !ok {
			return fmt.Errorf("Recipe %s doesn't exist", recipeName)
		}

		results := make([]string, 0)

		finished := false
		for finished != true {
			if current.InheritsExternal {
				if !excludeExternal {
					results = append(results, current.Inherits)
				}
				finished = true
			} else {
				current = rs[current.Inherits]
				results = append(results, current.Name)
			}
		}

		if reverse {
			results = utils.Reverse(results)
		}

		for _, result := range results {
			log.Println(result)
		}

		return nil
	},
}
