package recipes

import (
	gocontext "context"

	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/repository"
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:      "build",
	Usage:     "build a recipe",
	ArgsUsage: "<recipe>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "the tag to use when building this recipe",
			Value: "local",
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
			recipeName  = clicontext.Args().First()
			tag         = clicontext.String("tag")
			imagePrefix = clicontext.String("image-prefix")
		)

		recipe, err := recipes.GetRecipe(getRecipesDir(clicontext), recipeName)
		if err != nil {
			return err
		}

		s, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer s.Close()

		err = s.BuildRecipe(gocontext.Background(), recipe, tag, imagePrefix, nil)

		return err
	},
}
