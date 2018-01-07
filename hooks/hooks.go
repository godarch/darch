package hooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sort"

	"../stage"
	"../utils"
	"github.com/gobwas/glob"
)

// Hook A hook to be applied to images
type Hook struct {
	Name           string
	Path           string
	ExecutionOrder int
	IncludeImages  []string
	ExcludeImages  []string
}

type hookConfiguration struct {
	ExecutionOrder int
	IncludeImages  []string
	ExcludeImages  []string
}

type hookConfigurationJSON struct {
	ExecutionOrder *int      `json:"execution-order"`
	IncludeImages  *[]string `json:"include-images"`
	ExcludeImages  *[]string `json:"exclude-images"`
}

func buildDefaultHookEntry() hookConfiguration {
	return hookConfiguration{
		ExecutionOrder: 0,
		IncludeImages: []string{
			"*",
		},
		ExcludeImages: []string{},
	}
}

func getHooksConfiguration() (map[string]hookConfiguration, error) {

	if !utils.FileExists("/var/darch/hooks/hooks-config.json") {
		// No file exists, assume just a default entry.
		return map[string]hookConfiguration{
			"default": buildDefaultHookEntry(),
		}, nil
	}

	jsonData, err := ioutil.ReadFile("/var/darch/hooks/hooks-config.json")

	result := map[string]hookConfiguration{}
	jsonDesrialized := map[string]hookConfigurationJSON{}

	err = json.Unmarshal(jsonData, &jsonDesrialized)

	if err != nil {
		return result, err
	}

	// make sure there is a "_default" entry
	defaultEntry := buildDefaultHookEntry()
	if defaultEntrySerialized, ok := jsonDesrialized["_default"]; ok {
		if defaultEntrySerialized.ExecutionOrder != nil {
			defaultEntry.ExecutionOrder = *defaultEntrySerialized.ExecutionOrder
		}
		if defaultEntrySerialized.IncludeImages != nil {
			defaultEntry.IncludeImages = *defaultEntrySerialized.IncludeImages
		}
		if defaultEntrySerialized.ExcludeImages != nil {
			defaultEntry.ExcludeImages = *defaultEntrySerialized.ExcludeImages
		}
	}
	result["_default"] = defaultEntry

	// set default values on field that weren't set on hook configurtions.
	for k, v := range jsonDesrialized {
		if k == "_default" {
			continue
		}
		newEntry := defaultEntry
		if v.ExecutionOrder != nil {
			newEntry.ExecutionOrder = *v.ExecutionOrder
		}
		if v.IncludeImages != nil {
			newEntry.IncludeImages = *v.IncludeImages
		}
		if v.ExcludeImages != nil {
			newEntry.ExcludeImages = *v.ExcludeImages
		}
		result[k] = newEntry
	}

	return result, nil
}

// GetHook Get a hook by a name
func GetHook(name string) (Hook, error) {
	result := Hook{}

	if len(name) == 0 {
		return result, fmt.Errorf("A name is required")
	}

	result.Path = path.Join("/var/darch/hooks/" + name)

	if !utils.DirectoryExists(result.Path) {
		return result, fmt.Errorf("The hook %s doesn't exist", name)
	}

	configuration, err := getHooksConfiguration()

	if err != nil {
		return result, err
	}

	var config hookConfiguration
	if val, ok := configuration[name]; ok {
		config = val
	} else {
		config = configuration["_default"]
	}

	result.ExecutionOrder = config.ExecutionOrder
	result.IncludeImages = config.IncludeImages
	result.ExcludeImages = config.ExcludeImages

	return result, nil
}

// GetHooks Get all the available hooks
func GetHooks() ([]Hook, error) {
	result := make(map[string]Hook, 0)

	configuration, err := getHooksConfiguration()

	if err != nil {
		return nil, err
	}

	hooks, err := utils.GetChildDirectories("/var/darch/hooks")

	if err != nil {
		return nil, err
	}

	// This is used for sorting later
	executionOrderGroup := make(map[int][]string)

	for _, hook := range hooks {
		newHook := Hook{
			Name: hook,
			Path: path.Join("/var/darch/hooks", hook),
		}
		var config hookConfiguration
		if val, ok := configuration[newHook.Name]; ok {
			config = val
		} else {
			config = configuration["_default"]
		}
		newHook.ExecutionOrder = config.ExecutionOrder
		newHook.IncludeImages = config.IncludeImages
		newHook.ExcludeImages = config.ExcludeImages
		result[newHook.Name] = newHook

		if value, ok := executionOrderGroup[newHook.ExecutionOrder]; ok {
			executionOrderGroup[newHook.ExecutionOrder] = append(value, newHook.Name)
		} else {
			executionOrderGroup[newHook.ExecutionOrder] = []string{
				newHook.Name,
			}
		}
	}

	// Store the items in [ExecutionOrder->Name] order
	var executionOrderGroupKeys []int
	resultSorted := make([]Hook, 0)
	for k := range executionOrderGroup {
		executionOrderGroupKeys = append(executionOrderGroupKeys, k)
	}
	sort.Ints(executionOrderGroupKeys)
	for _, executionOrder := range executionOrderGroupKeys {
		executionOrderHooks := executionOrderGroup[executionOrder]
		sort.Strings(executionOrderHooks)
		for _, hook := range executionOrderHooks {
			resultSorted = append(resultSorted, result[hook])
		}
	}

	return resultSorted, nil
}

// AppliesToStagedTag Determines if a hook applies to the given staged tag
func AppliesToStagedTag(hook Hook, tag stage.StagedItemTag) bool {
	// First, let's see if we globbed the image
	for _, includeImage := range hook.IncludeImages {
		g := glob.MustCompile(includeImage)
		if g.Match(tag.FullName) {
			// This image has been included, but now, let's see if we excluded it
			for _, excludeImage := range hook.ExcludeImages {
				g = glob.MustCompile(excludeImage)
				if g.Match(tag.FullName) {
					// Someone doesn't want to apply this hook to this tag!
					return false
				}
			}
			// Nobody excluded us (after inclusion), so we are clear!
			return true
		}
	}
	return false
}
