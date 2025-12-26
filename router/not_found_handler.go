package router

import "net/http"

type NotFoundHandler struct{}

func (n NotFoundHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lg.Info().Str("method", req.Method).Str("path", req.URL.Path).Msg("Request not found")
	http.NotFound(w, req)
}
