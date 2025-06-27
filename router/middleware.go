package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware interface {
	// Handle processes the request. If an error occurs, it should return an error and the corresponding error code.
	Handle(res BaseResource, w http.ResponseWriter, req *http.Request, ps httprouter.Params) (error, int)
}
