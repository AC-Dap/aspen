package config

import (
	"aspen/router"
	"encoding/json"
	"fmt"
)

// We configure middleware as a map from type to MiddlewareConfig
type AllMiddlewareConfigs map[string]MiddlewareConfig

type MiddlewareConfig struct {
	Disabled bool
	Params   map[string]any
}

// Middleware define arbitrary parameters, so the best we can do is `any`.
type MiddlewareParams any
type MiddlewareContructor[P MiddlewareParams] func(P) router.Middleware

// A parser takes []byte JSON data and parses it using the relevant constructor into a middleware instance.
type MiddlewareParser func([]byte) (router.Middleware, error)

var globalMiddlewareParsersMap = make(map[string]MiddlewareParser)
var globalMiddlewareParamsMap = make(map[string]MiddlewareParams)

func RegisterMiddlewareConstructor[P MiddlewareParams](middlewareType string, constructor MiddlewareContructor[P]) error {
	// Check if this type alrady exists
	if _, ok := globalMiddlewareParsersMap[middlewareType]; ok {
		return fmt.Errorf("\"%s\" middleware constructor has already been registered", middlewareType)
	}

	// Save the parameters type for this middleware
	var params P
	globalMiddlewareParamsMap[middlewareType] = params

	// Create parser function
	parser := func(rawJson []byte) (router.Middleware, error) {
		var params P
		err := json.Unmarshal(rawJson, &params)
		if err != nil {
			return nil, fmt.Errorf("error parsing \"%s\" params: %w", middlewareType, err)
		}

		return constructor(params), nil
	}

	lg.Debug().Str("middleware", middlewareType).Msg("Registered middleware constructor")
	globalMiddlewareParsersMap[middlewareType] = parser
	return nil
}

// AvailableMiddleware returns a list of all registered middleware types.
func AvailableMiddleware() []string {
	var mTypes = make([]string, 0, len(globalMiddlewareParsersMap))
	for mType := range globalMiddlewareParsersMap {
		mTypes = append(mTypes, mType)
	}
	return mTypes
}

func (c AllMiddlewareConfigs) Parse() ([]router.Middleware, error) {
	// First verify that every registered middleware is present. We want to be explicit
	// with configurations.
	for mType := range globalMiddlewareParsersMap {
		if _, ok := c[mType]; !ok {
			return nil, fmt.Errorf("config missing for \"%s\" middleware", mType)
		}
	}

	var middlewares = make([]router.Middleware, 0, len(c))
	for mType, config := range c {
		parser, ok := globalMiddlewareParsersMap[mType]
		if !ok {
			return nil, fmt.Errorf("unable to find \"%s\" middleware constructor", mType)
		}

		if config.Disabled {
			continue
		}

		// Try parsing
		rawParams, err := json.Marshal(config.Params)
		if err != nil {
			return nil, fmt.Errorf("unable to read \"%s\" parameters", mType)
		}
		middleware, err := parser(rawParams)
		if err != nil {
			return nil, fmt.Errorf("unable to parse \"%s\" parameters", mType)
		}

		middlewares = append(middlewares, middleware)
	}

	return middlewares, nil
}
