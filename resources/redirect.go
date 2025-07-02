package resources

import (
	"aspen/router"
	"aspen/utils"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RedirectResource struct {
	host string
	path utils.Path
	router.BaseResource
}

type RedirectParams struct {
	Host string
	Path string
}

func NewRedirectResource(base router.BaseResource, params RedirectParams) router.Resource {
	return &RedirectResource{
		host:         params.Host,
		path:         utils.ParsePath(params.Path),
		BaseResource: base,
	}
}

func (rr *RedirectResource) Start() error {
	rr.BaseResource.Status = router.Started
	return nil
}

func (rr *RedirectResource) Stop() error {
	rr.BaseResource.Status = router.Stopped
	return nil
}

func (rr *RedirectResource) AddHandlers(path string, router *router.RouterInstance) error {
	// Check that the proxy path and given path have matching variables
	if !rr.path.IsProxyCompatible(utils.ParsePath(path)) {
		return fmt.Errorf("proxy path %s is not compatible with redirect path %s", path, rr.path)
	}

	router.GET(path, rr.BaseResource, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		constructedPath := rr.host + rr.path.ConstructPath(ps)

		// Redirect to the destination host and path
		http.Redirect(w, req, constructedPath, http.StatusFound)
	})

	return nil
}
