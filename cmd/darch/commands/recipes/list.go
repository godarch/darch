package recipes

import (
	"fmt"

	"github.com/pauldotknopf/darch/recipes"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list all recipes",
	Action: func(clicontext *cli.Context) error {
		rs, err := recipes.GetAllRecipes(getRecipesDir(clicontext))
		if err != nil {
			return err
		}

		for _, r := range rs {
			fmt.Println(r.Name)
		}

		return nil
	},
}
