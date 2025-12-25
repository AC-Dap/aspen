package middleware

import (
	"aspen/logging"
	"aspen/router"
	"aspen/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

var lg = logging.NewTaggedLogger("Logger")

type Logger struct{}

type LoggerParams struct{}

func NewLogger(params LoggerParams) router.Middleware {
	return &Logger{}
}

// Handle logs the request method, path, and the resource handling this request.
// It always continues the request after.
func (l *Logger) Handle(
	res router.BaseResource,
	trw utils.TrackingResponseWriter,
	req *http.Request,
	ps httprouter.Params,
	remainingMiddleware []router.MiddlewareHandler,
	requestHandler httprouter.Handle,
) {
	lg.Info().Str("method", req.Method).Str("path", req.URL.Path).Str("resource", res.GetID()).Msg("Request received")

	continueRequest(res, trw, req, ps, remainingMiddleware, requestHandler)

	var event *zerolog.Event
	if 500 <= trw.Status() && trw.Status() < 600 {
		// Internal server error, log as warning
		event = lg.Warn()
	} else {
		event = lg.Info()
	}
	event.
		Str("method", req.Method).
		Str("path", req.URL.Path).
		Str("resource", res.GetID()).
		Int("status", trw.Status()).
		Msg("Request finished")
}
