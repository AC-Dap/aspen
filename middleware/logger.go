package middleware

import (
	"aspen/logging"
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var lg = logging.NewTaggedLogger("Logger")

type Logger struct{}

type LoggerParams struct{}

func NewLogger(params LoggerParams) router.Middleware {
	return &Logger{}
}

// Handle logs the request method, path, and the resource handling this request.
// It always continues the request after.
func (l *Logger) Handle(res router.BaseResource, w http.ResponseWriter, req *http.Request, ps httprouter.Params) bool {
	lg.Info().Str("method", req.Method).Str("path", req.URL.Path).Str("resource", res.GetID()).Msg("Request received")
	return true
}
