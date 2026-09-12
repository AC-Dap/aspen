package config

import "fmt"

// We configure middleware as a map from type to MiddlewareConfig
type AllMiddlewareConfigs map[string]MiddlewareConfig

type MiddlewareConfig struct {
	Disabled bool
	Params   map[string]any
	Priority int
}

// Middleware define arbitrary parameters, so the best we can do is `any`.
type MiddlewareParams any

var globalMiddlewareParamsMap = make(map[string]MiddlewareParams)

// RegisterMiddlewareSchema saves the mapping from middlewareType to params
func RegisterMiddlewareSchema[P ResourceParams](middlewareType string) error {
	if _, ok := globalMiddlewareParamsMap[middlewareType]; ok {
		return fmt.Errorf("\"%s\" middleware schema has already been registered", middlewareType)
	}

	// Save the parameters type for this resource
	var params P
	globalMiddlewareParamsMap[middlewareType] = params

	lg.Debug().Str("type", middlewareType).Interface("params", params).Msg("Registered middleware schema")
	return nil
}

// AvailableMiddleware returns a list of all registered middleware types.
func AvailableMiddleware() []string {
	var names = make([]string, 0, len(globalMiddlewareParamsMap))
	for name := range globalMiddlewareParamsMap {
		names = append(names, name)
	}
	return names
}

// GetMiddlewareParams retrieves the parameters for a given middleware type.
func GetMiddlewareParams(middlewareType string) (ResourceParams, error) {
	params, ok := globalMiddlewareParamsMap[middlewareType]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" middleware parameters", middlewareType)
	}
	return params, nil
}
