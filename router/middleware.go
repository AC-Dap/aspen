package router

import (
	"aspen/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareHandler func(
	BaseResource,
	utils.TrackingResponseWriter,
	*http.Request,
	httprouter.Params,
	[]MiddlewareHandler,
	httprouter.Handle,
)

type Middleware interface {
	// Handle processes the request, and recursively calls the next middleware handler to continue
	// the request. If the middleware decides to cancel the request, it is responsible for writing
	// the response and returning early.
	Handle(
		res BaseResource,
		trw utils.TrackingResponseWriter,
		req *http.Request,
		ps httprouter.Params,
		remainingMiddleware []MiddlewareHandler,
		requestHandler httprouter.Handle,
	)
}
