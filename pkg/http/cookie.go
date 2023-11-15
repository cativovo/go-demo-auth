package http

import (
	"net/http"

	"github.com/cativovo/go-demo-auth/pkg/auth"
)

func createCookie(name string, value string, maxAge int) *http.Cookie {
	// https://www.alexedwards.net/blog/working-with-cookies-in-go
	return &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		cookie = createCookie(cookie.Name, cookie.Value, -1)
		http.SetCookie(w, cookie)
	}
}

func createTokenCookie(t auth.Token) (*http.Cookie, *http.Cookie) {
	accessTokenCookie := createCookie("access_token", t.AccessToken, t.ExpiresIn)
	refreshTokenCookie := createCookie("refresh_token", t.RefreshToken, 0)

	return accessTokenCookie, refreshTokenCookie
}
