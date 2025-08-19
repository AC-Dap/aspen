package resources

import (
	"aspen/auth"
	"aspen/router"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthResource struct {
	loginPage string
	router.BaseResource
}

type AuthParams struct {
	LoginPage string
}

func NewAuthResource(base router.BaseResource, params AuthParams) router.Resource {
	return &AuthResource{
		BaseResource: base,
		loginPage:    params.LoginPage,
	}
}

// AddHandlers adds auth-related pages and API routes.
func (a *AuthResource) AddHandlers(path string, router *router.RouterInstance) error {
	// Auth-related HTML pages
	router.GET(path+"/login", a.BaseResource, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		http.ServeFile(w, req, a.loginPage)
	})

	// Sign-up/Login/Logout API paths
	router.POST(path+"/create_user", a.BaseResource, create_user)
	router.POST(path+"/login", a.BaseResource, login)
	router.POST(path+"/logout", a.BaseResource, logout)
	return nil
}

func create_user(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	// Verify username and password are ok
	if body.Username == "" || body.Password == "" {
		http.Error(w, "Username and password must not be empty", http.StatusBadRequest)
		return
	}

	// Try adding to database
	err := auth.CreateUser(body.Username, body.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	// Try logging in with the provided credentials
	err := auth.VerifyUserCredentials(body.Username, body.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to login: %v", err), http.StatusUnauthorized)
		return
	}

	// Generate new access token
	token, err := auth.LoginUser(body.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate access token: %v", err), http.StatusInternalServerError)
	}

	auth.WriteAccessTokenCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Read access token cookie
	token, err := auth.ReadAccessTokenCookie(r)
	if err != nil {
		// Cookie not found, treat as logged out
		w.WriteHeader(http.StatusOK)
		return
	}

	// Try logging out
	err = auth.LogoutUser(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to logout: %v", err), http.StatusInternalServerError)
		return
	}

	auth.DeleteAccessTokenCookie(w)
	w.WriteHeader(http.StatusOK)
}
