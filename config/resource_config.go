package config

import "fmt"

type ResourceConfig struct {
	Type   string
	Params map[string]any
}

// Resources define arbitrary parameters, so the best we can do is `any`.
type ResourceParams any

var globalResourceParamsMap = make(map[string]ResourceParams)

// RegisterResourceSchema saves the mapping from resourceType to params
func RegisterResourceSchema[P ResourceParams](resourceType string) error {
	if _, ok := globalResourceParamsMap[resourceType]; ok {
		return fmt.Errorf("\"%s\" resource schema has already been registered", resourceType)
	}

	// Save the parameters type for this resource
	var params P
	globalResourceParamsMap[resourceType] = params

	lg.Debug().Str("type", resourceType).Interface("params", params).Msg("Registered resource schema")
	return nil
}

// AvailableResources returns a list of all registered resource types.
func AvailableResources() []string {
	var names = make([]string, 0, len(globalResourceParamsMap))
	for name := range globalResourceParamsMap {
		names = append(names, name)
	}
	return names
}

// GetResourceParams retrieves the parameters for a given resource type.
func GetResourceParams(resourceType string) (ResourceParams, error) {
	params, ok := globalResourceParamsMap[resourceType]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" resource parameters", resourceType)
	}
	return params, nil
}
