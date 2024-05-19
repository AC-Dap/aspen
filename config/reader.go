package config

import (
	. "dashboard/types"
	"github.com/pelletier/go-toml/v2"
	"os"
)

func Read() ResourcesConfig {
	configFile, err := os.ReadFile("resources.toml")
	if err != nil {
		panic(err)
	}

	var resources ResourcesConfig
	err = toml.Unmarshal(configFile, &resources)
	if err != nil {
		panic(err)
	}
	return resources
}

func Validate(resources ResourcesConfig) {
	// Verify individual resource fields
	for _, resource := range resources.Resources {
		prefix := "Error in resource " + resource.Name + ": "
		if resource.Name == "" {
			panic(prefix + "resource name cannot be empty")
		}
		if resource.Route == "" {
			panic(prefix + "resource route cannot be empty")
		}
		if resource.Source == "" {
			panic(prefix + "resource source cannot be empty")
		}
	}

	// Ensure no duplicate names
	names := make(map[string]bool)
	for _, resource := range resources.Resources {
		if names[resource.Name] {
			panic("Duplicate resource name: " + resource.Name)
		}
		names[resource.Name] = true
	}

	// We require a dashboard and auth resource
	if !names["dashboard"] {
		panic("Missing dashboard resource")
	}
	if !names["auth"] {
		panic("Missing auth resource")
	}
}
