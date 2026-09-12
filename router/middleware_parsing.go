package router

import (
	"aspen/config"
	"encoding/json"
	"fmt"
)

type MiddlewareContructor[P config.MiddlewareParams] func(P) Middleware

// A parser takes []byte JSON data and parses it using the relevant constructor into a middleware instance.
type MiddlewareParser func([]byte) (Middleware, error)

var globalMiddlewareParsersMap = make(map[string]MiddlewareParser)

func RegisterMiddlewareConstructor[P config.MiddlewareParams](middlewareType string, constructor MiddlewareContructor[P]) error {
	err := config.RegisterMiddlewareSchema[P](middlewareType)
	if err != nil {
		return fmt.Errorf("error registering middleware schema: %w", err)
	}

	// Create parser function
	parser := func(rawJson []byte) (Middleware, error) {
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

func ParseMiddleware(c config.AllMiddlewareConfigs) ([]Middleware, error) {
	// First verify that every registered middleware is present. We want to be explicit
	// with configurations.
	for _, mType := range config.AvailableMiddleware() {
		if _, ok := c[mType]; !ok {
			return nil, fmt.Errorf("config missing for \"%s\" middleware", mType)
		}
	}

	var middlewares = make([]Middleware, 0, len(c))
	var priorities = make([]int, 0, len(c))
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
		priorities = append(priorities, config.Priority)
	}

	// Sort middlewares by their priorities
	for i := 0; i < len(middlewares); i++ {
		for j := i + 1; j < len(middlewares); j++ {
			if priorities[i] == priorities[j] {
				return nil, fmt.Errorf("found duplicate priorities %d", priorities[i])
			}
			if priorities[i] > priorities[j] {
				middlewares[i], middlewares[j] = middlewares[j], middlewares[i]
				priorities[i], priorities[j] = priorities[j], priorities[i]
			}
		}
	}

	return middlewares, nil
}
