package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/cativovo/go-demo-auth/pkg/user"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerAuthRoutes() {
	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.handleRegister)
		r.Post("/login", s.handleLogin)
		r.Post("/logout", s.handleLogout)
	})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	userCredentials := user.Credentials{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		Name:     r.PostFormValue("name"),
	}

	token, errors := s.userService.Register(userCredentials)
	if errors != nil {
		// TODO: handle each errors
		errorAlertTmpl.Execute(w, map[string]any{
			"Message": "Something went wrong...",
		})
		return
	}

	accessTokenCookie, refreshTokenCookie := createTokenCookie(token)

	w.Header().Add("HX-Trigger", "redirect-to-account")
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	token, err := s.authService.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			errorAlertTmpl.Execute(w, map[string]any{
				"Message": "Invalid username/password",
			})
			return
		default:
			errorAlertTmpl.Execute(w, map[string]any{
				"Message": "Something went wrong",
			})
			return
		}
	}

	accessTokenCookie, refreshTokenCookie := createTokenCookie(token)

	w.Header().Add("HX-Trigger", "redirect-to-account")
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	jwt, err := r.Cookie("jwt")
	if err != nil {
		return
	}

	if err := s.authService.Logout(jwt.Value); err != nil {
		fmt.Fprintln(w, "ooo")
		return
	}

	// clear cookie
	// redirect to login page
}

func createTokenCookie(t auth.Token) (http.Cookie, http.Cookie) {
	// https://www.alexedwards.net/blog/working-with-cookies-in-go
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    t.AccessToken,
		MaxAge:   t.ExpiresIn,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    t.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	return accessTokenCookie, refreshTokenCookie
}
