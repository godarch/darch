package recipes

import (
	"fmt"

	"github.com/pauldotknopf/darch/utils"
)

// Recipe A struct representing a recipe to be built.
type Recipe struct {
	Name             string
	RecipeDir        string
	RecipesDir       string
	Inherits         string
	InheritsExternal bool
}

func verifyDependencies(recipe Recipe, recipes map[string]Recipe, currentStack map[string]bool) error {
	if currentStack == nil {
		currentStack = make(map[string]bool, 0)
	}

	if recipe.InheritsExternal {
		// we reached the end, all good!
		return nil
	}

	if _, ok := currentStack[recipe.Inherits]; ok {
		// Cyclical dependency detected!
		return fmt.Errorf("Recipe %s has a cyclical dependency", recipe.Name)
	}

	// Make this image as having been traversed.
	currentStack[recipe.Name] = true

	if parent, ok := recipes[recipe.Inherits]; ok {
		return verifyDependencies(parent, recipes, currentStack)
	}

	return fmt.Errorf("Recipe defintion %s inherits from %s, which doesn't exist", recipe.Name, recipe.Inherits)
}

// GetAllRecipes Return all the recipes in a recipe directory
func GetAllRecipes(recipesDir string) (map[string]Recipe, error) {
	if len(recipesDir) == 0 {
		return nil, fmt.Errorf("An image directory must be provided")
	}

	recipeNames, err := utils.GetChildDirectories(recipesDir)

	if err != nil {
		return nil, err
	}

	recipes := make(map[string]Recipe, 0)

	for _, recipeName := range recipeNames {
		recipe, err := parseRecipe(recipesDir, recipeName)
		if err != nil {
			return nil, err
		}
		recipes[recipeName] = recipe
	}

	// verify dependencies are satisfied and no circular dependencies
	for _, recipe := range recipes {
		err := verifyDependencies(recipe, recipes, nil)
		if err != nil {
			return nil, err
		}
	}

	return recipes, nil
}

// GetRecipe Get a single recipe by name
func GetRecipe(recipesDir string, recipeName string) (Recipe, error) {
	allRecipes, err := GetAllRecipes(recipesDir)
	if err != nil {
		return Recipe{}, err
	}

	current, ok := allRecipes[recipeName]
	if !ok {
		return Recipe{}, fmt.Errorf("Recipe %s doesn't exist", recipeName)
	}

	return current, nil
}
