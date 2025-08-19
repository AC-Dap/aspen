package config

import (
	"aspen/router"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type MiddlewareConfig struct {
	Type   string
	Params map[string]any
}

// Middleware define arbitrary parameters, so the best we can do is `any`.
type MiddlewareParams = any
type MiddlewareContructor[P MiddlewareParams] = func(P) router.Middleware

// A parser takes []byte JSON data and parses it using the relevant constructor into a middleware instance.
type MiddlewareParser = func([]byte) (router.Middleware, error)
type MiddlewareParsers = map[string]MiddlewareParser

var globalMiddlewareMap = make(MiddlewareParsers)
var globalMiddlewareParamsMap = make(map[string]MiddlewareParams)

func RegisterMiddlewareConstructor[P MiddlewareParams](middlewareType string, constructor MiddlewareContructor[P]) error {
	// Check if this type alrady exists
	if _, ok := globalMiddlewareMap[middlewareType]; ok {
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

	log.Debug().Str("middleware", middlewareType).Msg("Registered middleware constructor")
	globalMiddlewareMap[middlewareType] = parser
	return nil
}

// AvailableMiddleware returns a list of all registered middleware names.
func AvailableMiddleware() []string {
	var names = make([]string, 0, len(globalMiddlewareMap))
	for name := range globalMiddlewareMap {
		names = append(names, name)
	}
	return names
}

func (m MiddlewareConfig) Parse() (router.Middleware, error) {
	parser, ok := globalMiddlewareMap[m.Type]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" middleware constructor", m.Type)
	}

	// Try parsing
	rawParams, err := json.Marshal(m.Params)
	if err != nil {
		return nil, fmt.Errorf("unable to read \"%s\" parameters", m.Type)
	}
	middleware, err := parser(rawParams)
	if err != nil {
		return nil, fmt.Errorf("unable to parse \"%s\" parameters", m.Type)
	}

	return middleware, nil
}
