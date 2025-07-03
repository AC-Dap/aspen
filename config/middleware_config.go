package config

import (
	"aspen/router"
	"fmt"
)

type MiddlewareConfig string

var globalMiddlewareMap = make(map[string]router.Middleware)

func RegisterMiddleware(name string, middleware router.Middleware) error {
	// Check if this type alrady exists
	if _, ok := globalMiddlewareMap[name]; ok {
		return fmt.Errorf("\"%s\" resource constructor has already been registered", name)
	}

	globalMiddlewareMap[name] = middleware
	return nil
}

func (m MiddlewareConfig) Parse() (router.Middleware, error) {
	middleware, ok := globalMiddlewareMap[string(m)]
	if !ok {
		return nil, fmt.Errorf("unable to find \"%s\" middleware", m)
	}
	return middleware, nil
}
