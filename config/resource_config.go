package config

import (
	"aspen/router"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ResourceConfig struct {
	ResourceType string
	Params       map[string]any
}

// Resources define arbitrary parameters, so the best we can do is `any`.
type ResourceParams = any
type ResourceContructor[P ResourceParams] = func(router.BaseResource, P) router.Resource

// A parser takes []byte JSON data and parses it using the relevant constructor into a resource instance.
type ResourceParser = func(router.BaseResource, []byte) (router.Resource, error)
type ResourceParsers = map[string]ResourceParser

var globalResourceMap = make(ResourceParsers)
var globalResourceParamsMap = make(map[string]ResourceParams)

func RegisterResourceConstructor[P ResourceParams](resourceType string, constructor ResourceContructor[P]) error {
	// Check if this type alrady exists
	if _, ok := globalResourceMap[resourceType]; ok {
		return fmt.Errorf("\"%s\" resource constructor has already been registered", resourceType)
	}

	// Save the parameters type for this resource
	var params P
	globalResourceParamsMap[resourceType] = params

	// Create parser function
	parser := func(base router.BaseResource, rawJson []byte) (router.Resource, error) {
		var params P
		err := json.Unmarshal(rawJson, &params)
		if err != nil {
			return nil, fmt.Errorf("error parsing \"%s\" params: %w", resourceType, err)
		}

		return constructor(base, params), nil
	}

	log.Debug().Str("resource", resourceType).Msg("Registered resource constructor")
	globalResourceMap[resourceType] = parser
	return nil
}

// AvailableResources returns a list of all registered resource names.
func AvailableResources() []string {
	var names = make([]string, 0, len(globalResourceMap))
	for name := range globalResourceMap {
		names = append(names, name)
	}
	return names
}

// GetResourceParams retrieves the parameters type for a given resource type.
func GetResourceParams(resourceType string) (ResourceParams, error) {
	params, ok := globalResourceParamsMap[resourceType]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" resource parameters", resourceType)
	}
	return params, nil
}

func (rc ResourceConfig) Parse(base router.BaseResource) (router.Resource, error) {
	parser, ok := globalResourceMap[rc.ResourceType]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" resource constructor", rc.ResourceType)
	}

	// Try parsing
	rawParams, err := json.Marshal(rc.Params)
	if err != nil {
		return nil, fmt.Errorf("unable to read \"%s\" parameters", rc.ResourceType)
	}
	newResource, err := parser(base, rawParams)
	if err != nil {
		return nil, fmt.Errorf("unable to parse \"%s\" parameters", rc.ResourceType)
	}

	return newResource, nil
}
