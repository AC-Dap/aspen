package auth

import (
	"net/http"
)

const AccessTokenCookie = "aspen_auth"

func ReadAccessTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookie)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func WriteAccessTokenCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(TokenLifespan.Seconds()),
	}

	http.SetCookie(w, cookie)
}

func DeleteAccessTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}
