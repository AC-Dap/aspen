package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func Read() ([]Resource, error) {
	configFile, err := os.ReadFile("resources.toml")
	if err != nil {
		return nil, err
	}

	var resources ConfigTOML
	err = toml.Unmarshal(configFile, &resources)
	if err != nil {
		return nil, err
	}
	return resources.Resources, nil
}

func Validate(resources []Resource) {
	// Verify individual resource fields
	for _, resource := range resources {
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
	for _, resource := range resources {
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
