package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware interface {
	// Handle processes the request, and returns a boolean deciding whether to continue the request handling.
	// Handle is responsible for writing the response if it decides to stop the request handling.
	Handle(res BaseResource, w http.ResponseWriter, req *http.Request, ps httprouter.Params) bool
}
