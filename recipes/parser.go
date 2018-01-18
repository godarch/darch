package recipes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pauldotknopf/darch/utils"
)

type recipeConfiguration struct {
	Inherits string `json:"inherits"`
}

func parseRecipe(recipesDir string, recipeName string) (Recipe, error) {
	recipe := Recipe{}

	if len(recipesDir) == 0 {
		return recipe, fmt.Errorf("A recipe directory must be provided")
	}

	if len(recipeName) == 0 {
		return recipe, fmt.Errorf("A recipe name must be provided")
	}

	recipe.RecipesDir = utils.ExpandPath(recipesDir)
	recipe.RecipeDir = path.Join(recipe.RecipesDir, recipeName)
	recipe.Name = recipeName

	if !utils.DirectoryExists(recipe.RecipeDir) {
		return recipe, fmt.Errorf("Image directory %s doesn't exist", recipe.RecipeDir)
	}

	recipeConfiguration, err := loadRecipeConfiguration(recipe)

	if err != nil {
		return recipe, err
	}

	if strings.HasPrefix(recipeConfiguration.Inherits, "external:") {
		recipe.InheritsExternal = true
		recipe.Inherits = recipeConfiguration.Inherits[len("external:"):len(recipeConfiguration.Inherits)]
	} else {
		recipe.InheritsExternal = false
		recipe.Inherits = recipeConfiguration.Inherits
	}

	return recipe, nil
}

func loadRecipeConfiguration(recipe Recipe) (recipeConfiguration, error) {
	imageConfigurationPath := path.Join(image.ImageDir, "config.json")
	imageConfiguration := imageConfiguration{}

	if !utils.FileExists(imageConfigurationPath) {
		return nil, fmt.Errorf("No configuration file exists at %s", imageConfigurationPath)
	}

	jsonData, err := ioutil.ReadFile(imageConfigurationPath)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &imageConfiguration)

	if err != nil {
		return nil, err
	}

	if len(imageConfiguration.Inherits) == 0 {
		return nil, fmt.Errorf("No inherit property given for image %s", image.Name)
	}

	return &imageConfiguration, nil
}
