package middleware

import (
	"aspen/router"
	"aspen/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func continueRequest(
	res router.BaseResource,
	trw utils.TrackingResponseWriter,
	req *http.Request,
	ps httprouter.Params,
	remainingMiddleware []router.MiddlewareHandler,
	requestHandler httprouter.Handle,
) {
	if len(remainingMiddleware) == 0 {
		requestHandler(trw, req, ps)
	} else {
		remainingMiddleware[0](res, trw, req, ps, remainingMiddleware[1:], requestHandler)
	}
}
