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
		r.Get("/logout", s.handleLogout)
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

	w.Header().Add("HX-Location", "/")
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
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

	w.Header().Add("HX-Location", "/")
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if token, getCookieErr := r.Cookie("access_token"); getCookieErr == nil {
		if logoutErr := s.authService.Logout(token.Value); logoutErr != nil {
			fmt.Fprintln(w, logoutErr)
		}
	}

	w.Header().Add("Location", "/auth-page/login")
	clearCookie(w, r)
	w.WriteHeader(http.StatusFound)
}
