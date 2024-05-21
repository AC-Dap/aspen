package auth

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request, users []User) {
	fmt.Fprintf(w, "Login page\n%v", users)
}

func verifyHandler(w http.ResponseWriter, r *http.Request, users []User) {
	fmt.Fprintf(w, "Auth API\n%v", users)
}

func logoutHandler(w http.ResponseWriter, r *http.Request, users []User) {
	fmt.Fprintf(w, "Logout page\n%v", users)
}

func handleWithUsers(handler func(http.ResponseWriter, *http.Request, []User), users []User) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, users)
	}
}

// AddRoutes
// Adds the necessary routes for our auth API, starting from the given base route
func AddRoutes(baseRoute string, serveMux *http.ServeMux, users []User) {
	serveMux.HandleFunc(fmt.Sprintf("GET %s/login", baseRoute),
		handleWithUsers(loginHandler, users))
	serveMux.HandleFunc(fmt.Sprintf("POST %s/verify", baseRoute),
		handleWithUsers(verifyHandler, users))
	serveMux.HandleFunc(fmt.Sprintf("POST %s/logout", baseRoute),
		handleWithUsers(logoutHandler, users))
}
