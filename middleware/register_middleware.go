package middleware

import "aspen/config"

func RegisterMiddleware() {
	config.RegisterMiddlewareConstructor[LoggerParams]("logger", NewLogger)
	config.RegisterMiddlewareConstructor[AuthParams]("auth", NewAuth)
}
