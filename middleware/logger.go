package middleware

import (
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

// Handle logs the request method, path, and the resource handling this request.
func (l Logger) Handle(res router.BaseResource, w http.ResponseWriter, req *http.Request, ps httprouter.Params) (error, int) {
	log.Info().Str("method", req.Method).Str("path", req.URL.Path).Str("resource", res.GetID()).Msg("Request received")
	return nil, http.StatusOK
}
