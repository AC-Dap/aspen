package resources

import (
	"aspen/router"
	"aspen/utils"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ProxyResource struct {
	host    string
	path    utils.Path
	methods []string
	router.BaseResource
}

type ProxyParams struct {
	Host    string
	Path    string
	Methods []string
}

func NewProxyResource(base router.BaseResource, params ProxyParams) router.Resource {
	return &ProxyResource{
		host:         params.Host,
		path:         utils.ParsePath(params.Path),
		methods:      params.Methods,
		BaseResource: base,
	}
}

func (pr *ProxyResource) Start() error {
	pr.BaseResource.Status = router.Started
	return nil
}

func (pr *ProxyResource) Stop() error {
	pr.BaseResource.Status = router.Stopped
	return nil
}

func (pr *ProxyResource) AddHandlers(path string, router *router.RouterInstance) error {
	// Check that the proxy path and given path have matching variables
	if !pr.path.IsProxyCompatible(utils.ParsePath(path)) {
		return fmt.Errorf("proxy path %s is not compatible with redirect path %s", path, pr.path)
	}

	// Every proxy request uses the same client (safe for concurrent use)
	proxyClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// Register the proxy handler for the specified methods
	for _, method := range pr.methods {
		router.Handle(method, path, pr.BaseResource, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			constructedPath := pr.host + pr.path.ConstructPath(ps)

			// Create a new request to the destination host and path
			proxyReq, err := http.NewRequest(method, constructedPath, req.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating proxy request: %v", err), http.StatusInternalServerError)
				return
			}
			proxyReq.Header = req.Header

			// Forward the request to the destination host
			resp, err := proxyClient.Do(proxyReq)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error forwarding request: %v", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			// Relay response back to the original client
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		})
	}

	return nil
}
