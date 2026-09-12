package middleware

import "aspen/router"

func RegisterMiddleware() {
	router.RegisterMiddlewareConstructor("logger", NewLogger)
	router.RegisterMiddlewareConstructor("auth", NewAuth)
}
