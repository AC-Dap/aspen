package middleware

import (
	"aspen/auth"
	"aspen/router"
	"aspen/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	// Path under which the auth resource is served
	Path string
}

type AuthParams struct {
	Path string
}

func NewAuth(params AuthParams) router.Middleware {
	return &Auth{
		Path: params.Path,
	}
}

func (a *Auth) Handle(
	res router.BaseResource,
	trw utils.TrackingResponseWriter,
	req *http.Request,
	ps httprouter.Params,
	remainingMiddleware []router.MiddlewareHandler,
	requestHandler httprouter.Handle,
) {
	// Check if this route is protected
	if len(res.AccessRoles) == 0 {
		// No roles required, allow access
		continueRequest(res, trw, req, ps, remainingMiddleware, requestHandler)
		return
	}

	token, err := auth.ReadAccessTokenCookie(req)
	if err == nil {
		// Validate the token has access to at least one role
		for _, role := range res.AccessRoles {
			valid, err := auth.VerifyAccessToken(token, role)
			if err == nil && valid {
				// Token is valid for this role, allow access
				continueRequest(res, trw, req, ps, remainingMiddleware, requestHandler)
				return
			}
		}
	}

	// If we reach here, the token is either missing or invalid for all roles
	http.Redirect(trw, req, a.Path+"/login", http.StatusFound)
}
