package recipes

import (
	"fmt"
	"log"
	"sort"

	"github.com/pauldotknopf/darch/pkg/recipes"
	"github.com/urfave/cli"
)

var childrenCommand = cli.Command{
	Name:  "children",
	Usage: "list all the children for a recipe",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "reverse",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			recipeName = clicontext.Args().First()
			reverse    = clicontext.Bool("reverse")
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

		for _, r := range rs {
			if r.Inherits == current.Name {
				results = append(results, r.Name)
			}
		}

		if reverse {
			sort.Sort(sort.Reverse(sort.StringSlice(results)))
		}

		for _, result := range results {
			log.Println(result)
		}

		return nil
	},
}
