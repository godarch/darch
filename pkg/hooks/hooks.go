package hooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"

	"github.com/gobwas/glob"
	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/utils"
)

var (
	// DefaultHooksConfigLocation Where the hooks configuration lives.
	DefaultHooksConfigLocation = "/var/lib/darch/hooks/hooks-config.json"
	// DefaultHooksPath Where the hooks live.
	DefaultHooksPath = "/var/lib/darch/hooks/"
)

// Hook A hook to be applied to images
type Hook struct {
	Name           string
	Path           string
	HooksPath      string
	ExecutionOrder int
	NameWithOrder  string //ExecutionOrder_Name
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

	if !utils.FileExists(DefaultHooksConfigLocation) {
		// No file exists, assume just a default entry.
		return map[string]hookConfiguration{
			"_default": buildDefaultHookEntry(),
		}, nil
	}

	jsonData, err := ioutil.ReadFile(DefaultHooksConfigLocation)

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
		return result, fmt.Errorf("a name is required")
	}

	result.Name = name
	result.HooksPath = DefaultHooksPath
	result.Path = path.Join(result.HooksPath, name)

	if !utils.DirectoryExists(result.Path) {
		return result, fmt.Errorf("the hook %s doesn't exist", name)
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
	result.NameWithOrder = fmt.Sprintf("%08d_%s", result.ExecutionOrder, result.Name)

	return result, nil
}

// GetHooks Get all the available hooks
func GetHooks() ([]Hook, error) {
	result := make(map[string]Hook, 0)

	configuration, err := getHooksConfiguration()

	if err != nil {
		return nil, err
	}

	if !utils.DirectoryExists(DefaultHooksPath) {
		return nil, nil
	}

	hooks, err := utils.GetChildDirectories(DefaultHooksPath)

	if err != nil {
		return nil, err
	}

	// This is used for sorting later
	executionOrderGroup := make(map[int][]string)

	for _, hook := range hooks {
		newHook := Hook{
			Name:      hook,
			HooksPath: DefaultHooksPath,
		}
		newHook.Path = path.Join(newHook.HooksPath, newHook.Name)
		var config hookConfiguration
		if val, ok := configuration[newHook.Name]; ok {
			config = val
		} else {
			config = configuration["_default"]
		}
		newHook.ExecutionOrder = config.ExecutionOrder
		newHook.IncludeImages = config.IncludeImages
		newHook.ExcludeImages = config.ExcludeImages
		newHook.NameWithOrder = fmt.Sprintf("%08d_%s", newHook.ExecutionOrder, newHook.Name)
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

// AppliesToImage Determines if a hook applies to the given image.
func AppliesToImage(hook Hook, imageRef reference.ImageRef) bool {
	// First, let's see if we globbed the image
	for _, includeImage := range hook.IncludeImages {
		g := glob.MustCompile(includeImage)
		if g.Match(imageRef.FullName()) {
			// This image has been included, but now, let's see if we excluded it
			for _, excludeImage := range hook.ExcludeImages {
				g = glob.MustCompile(excludeImage)
				if g.Match(imageRef.FullName()) {
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

// PrintHookHelp Print the help for a hook
func PrintHookHelp(hook Hook) error {
	hookFile := path.Join(hook.Path, "hook")
	if !utils.FileExists(hookFile) {
		return fmt.Errorf("Hook script %s doesn't exist", hookFile)
	}

	cmd := exec.Command("/bin/bash", "-c", ". "+hookFile+" && help")
	cmd.Dir = hook.HooksPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
