package recipes

import (
	"github.com/disiqueira/gotree"
	"github.com/pauldotknopf/darch/recipes"
	"github.com/pauldotknopf/darch/utils"
	"github.com/urfave/cli"
)

var treeCommand = cli.Command{
	Name:  "tree",
	Usage: "list all recipes in a tree",
	Action: func(clicontext *cli.Context) error {
		rs, err := recipes.GetAllRecipes(getRecipesDir(clicontext))
		if err != nil {
			return err
		}

		externalImages := make([]string, 0)

		for _, r := range rs {
			if r.InheritsExternal {
				externalImages = append(externalImages, r.Inherits)
			}
		}

		// this will be our root items
		externalImages = utils.RemoveDuplicates(externalImages)

		var rootNode gotree.GTStructure

		for _, externalImage := range externalImages {
			var externalImageNode gotree.GTStructure
			externalImageNode.Name = externalImage
			for _, r := range rs {
				if r.InheritsExternal && r.Inherits == externalImage {
					var childNode gotree.GTStructure
					childNode.Name = r.Name
					for _, child := range buildTreeRecursively(r, rs) {
						childNode.Items = append(childNode.Items, child)
					}
					externalImageNode.Items = append(externalImageNode.Items, childNode)
				}
			}
			rootNode.Items = append(rootNode.Items, externalImageNode)
		}

		gotree.PrintTree(rootNode)

		return nil
	},
}

func buildTreeRecursively(parentDefinition recipes.Recipe, rs map[string]recipes.Recipe) []gotree.GTStructure {
	children := make([]gotree.GTStructure, 0)

	for _, childRecipeDefinition := range rs {
		if childRecipeDefinition.Inherits == parentDefinition.Name {
			var childNode gotree.GTStructure
			childNode.Name = childRecipeDefinition.Name

			for _, child := range buildTreeRecursively(childRecipeDefinition, rs) {
				childNode.Items = append(childNode.Items, child)
			}
			children = append(children, childNode)
		}
	}

	return children
}
