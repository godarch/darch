package recipes

import (
	"context"
	"fmt"

	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/repository"
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:      "build",
	Usage:     "build a recipe(s)",
	ArgsUsage: "<recipes>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "the tag to use when building this recipe",
			Value: "latest",
		},
		cli.StringFlag{
			Name:  "image-prefix, p",
			Usage: "the value to prepend to all image names (inherited and built)",
			Value: "",
		},
		cli.StringSliceFlag{
			Name: "environment, e",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			tag         = clicontext.String("tag")
			imagePrefix = clicontext.String("image-prefix")
			recipeNames = clicontext.Args()
			env         = clicontext.StringSlice("env")
		)

		if len(recipeNames) == 0 {
			return fmt.Errorf("no recipes provided")
		}

		allRecipes, err := recipes.GetAllRecipes(getRecipesDir(clicontext))
		if err != nil {
			return err
		}

		// First, let's make sure all the recipes we are building exist.
		for _, recipeName := range recipeNames {
			if _, ok := allRecipes[recipeName]; !ok {
				return fmt.Errorf("recipe %s doesn't exist", recipeName)
			}
		}

		session, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}

		// Now, let's go through each recipe and build it.
		for _, recipeName := range recipeNames {
			fmt.Printf("building %s...\n", recipeName)
			err := session.BuildRecipe(context.Background(), allRecipes[recipeName], tag, imagePrefix, env)
			if err != nil {
				return err
			}
			fmt.Printf("built %s\n", recipeName)
		}

		return err
	},
}
