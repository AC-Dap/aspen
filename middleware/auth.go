package middleware

import (
	"aspen/auth"
	"aspen/router"
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

func (a *Auth) Handle(res router.BaseResource, w http.ResponseWriter, req *http.Request, ps httprouter.Params) bool {
	// Check if this route is protected
	if len(res.AccessRoles) == 0 {
		// No roles required, allow access
		return true
	}

	token, err := auth.ReadAccessTokenCookie(req)
	if err == nil {
		// Validate the token has access to at least one role
		for _, role := range res.AccessRoles {
			valid, err := auth.VerifyAccessToken(token, role)
			if err == nil && valid {
				// Token is valid for this role, allow access
				return true
			}
		}
	}

	// If we reach here, the token is either missing or invalid for all roles
	http.Redirect(w, req, a.Path+"/login", http.StatusFound)
	return false
}
