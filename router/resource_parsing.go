package router

import (
	"aspen/config"
	"encoding/json"
	"fmt"
)

type ResourceContructor[P config.ResourceParams] func(BaseResource, P) Resource

// A parser takes []byte JSON data and parses it using the relevant constructor into a resource instance.
type ResourceParser func(BaseResource, []byte) (Resource, error)

var globalResourceParsersMap = make(map[string]ResourceParser)

func RegisterResourceConstructor[P config.ResourceParams](resourceType string, constructor ResourceContructor[P]) error {
	err := config.RegisterResourceSchema[P](resourceType)
	if err != nil {
		return fmt.Errorf("error registering resource schema: %w", err)
	}

	// Create parser function
	parser := func(base BaseResource, rawJson []byte) (Resource, error) {
		var params P
		err := json.Unmarshal(rawJson, &params)
		if err != nil {
			return nil, fmt.Errorf("error parsing \"%s\" params: %w", resourceType, err)
		}

		return constructor(base, params), nil
	}

	lg.Debug().Str("resource", resourceType).Msg("Registered resource constructor")
	globalResourceParsersMap[resourceType] = parser
	return nil
}

func ParseResource(base BaseResource, rc config.ResourceConfig) (Resource, error) {
	parser, ok := globalResourceParsersMap[rc.Type]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" resource constructor", rc.Type)
	}

	// Try parsing
	rawParams, err := json.Marshal(rc.Params)
	if err != nil {
		return nil, fmt.Errorf("unable to read \"%s\" parameters", rc.Type)
	}
	newResource, err := parser(base, rawParams)
	if err != nil {
		return nil, fmt.Errorf("unable to parse \"%s\" parameters", rc.Type)
	}

	return newResource, nil
}
