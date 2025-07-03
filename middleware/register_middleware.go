package middleware

import "aspen/config"

func RegisterMiddleware() {
	config.RegisterMiddleware("logger", Logger{})
}
