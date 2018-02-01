package recipes

import (
	"fmt"

	"github.com/godarch/darch/pkg/recipes"
	"github.com/godarch/darch/pkg/utils"
	"github.com/urfave/cli"
)

var builddepCommand = cli.Command{
	Name:      "build-dep",
	Usage:     "list dependencies for the given recipes",
	ArgsUsage: "<recipes>*N",
	Action: func(clicontext *cli.Context) error {
		var (
			recipeNames = clicontext.Args()
		)

		allRecipes, err := recipes.GetAllRecipes(getRecipesDir(clicontext))
		if err != nil {
			return err
		}

		if len(recipeNames) == 0 {
			// We want dependencies for all images.
			for _, recipe := range allRecipes {
				recipeNames = append(recipeNames, recipe.Name)
			}
		}

		// First, let's make sure all the recipes we are building exist.
		for _, recipeName := range recipeNames {
			if _, ok := allRecipes[recipeName]; !ok {
				return fmt.Errorf("recipe %s doesn't exist", recipeName)
			}
		}

		dependencies := make([]string, 0)

		for _, r := range allRecipes {
			if utils.Contains(recipeNames, r.Name) {
				parents := walkRecipeRecursively(r, allRecipes)
				parents = utils.Reverse(parents)
				dependencies = append(dependencies, parents...)
			}
		}

		dependencies = utils.RemoveDuplicates(dependencies)

		for _, dependency := range dependencies {
			fmt.Println(dependency)
		}

		return err
	},
}

func walkRecipeRecursively(r recipes.Recipe, rs map[string]recipes.Recipe) []string {
	result := make([]string, 0)
	result = append(result, r.Name)
	if !r.InheritsExternal {
		children := walkRecipeRecursively(rs[r.Inherits], rs)
		result = append(result, children...)
	}
	return result
}
