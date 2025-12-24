package middleware

import "aspen/config"

func RegisterMiddleware() {
	config.RegisterMiddlewareConstructor("logger", NewLogger)
	config.RegisterMiddlewareConstructor("auth", NewAuth)
}
