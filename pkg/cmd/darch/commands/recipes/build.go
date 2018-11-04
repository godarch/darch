package recipes

import (
	"context"
	"fmt"
	"github.com/godarch/darch/pkg/recipes"
	"github.com/godarch/darch/pkg/repository"
	"github.com/urfave/cli"
	"strings"
)

var buildCommand = cli.Command{
	Name:      "build",
	Usage:     "build a recipe(s)",
	ArgsUsage: "<recipes>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "tags, t",
			Usage: "the tag(s) to use when building the recipe",
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
			tags        = clicontext.String("tags")
			imagePrefix = clicontext.String("image-prefix")
			recipeNames = clicontext.Args()
			env         = clicontext.StringSlice("environment")
		)

		if len(recipeNames) == 0 {
			return fmt.Errorf("no recipes provided")
		}

		defaultTag, additionalTags, err := parseTags(tags)
		if err != nil {
			return err
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
			image, err := session.BuildRecipe(context.Background(), allRecipes[recipeName], defaultTag, imagePrefix, env)
			if err != nil {
				return err
			}
			fmt.Printf("built %s as %s\n", recipeName, image.FullName())
			// Add additional tags.
			if len(additionalTags) > 0 {
				for _, tag := range additionalTags {
					newImageRef, err := image.WithTag(tag)
					if err != nil {
						return err
					}
					fmt.Printf("tagging as %s\n", newImageRef.FullName())
					err = session.TagImage(context.Background(), image, newImageRef)
					if err != nil {
						return err
					}
				}
			}
		}

		return err
	},
}

func parseTags(tags string) (string, []string, error) {
	if len(tags) == 0 {
		return "", nil, fmt.Errorf("invalid tag")
	}

	split := strings.Split(tags, ",")

	for _, tag := range split {
		if len(tag) == 0 {
			return "", nil, fmt.Errorf("invalid tag")
		}
	}

	return split[0], split[1:], nil
}
