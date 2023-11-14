package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/cativovo/go-demo-auth/pkg/user"
)

func (s *Server) registerAuthRoutes() {
	s.router.Post("/register", s.handleRegister)
	s.router.Post("/login", s.handleLogin)
	s.router.Post("/logout", s.handleLogout)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	userCredentials := user.Credentials{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		Name:     r.PostFormValue("name"),
	}

	// token, err
	_, err := s.userService.Register(userCredentials)
	if err != nil {
		switch {
		// case errors.Is(err, user.ErrEmailAlreadyUsed):
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	fmt.Fprintln(w, "Email is already used")
		// case errors.Is(err, user.ErrPasswordTooShort):
		// 	fmt.Fprintln(w, "Password is too short")
		// case errors.Is(err, user.ErrFieldRequired):
		// 	message := fmt.Sprintf("%s", err)
		// 	fmt.Fprintln(w, message)
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "oops something went wrong")
		}

		return
	}

	fmt.Fprintln(w, "registered")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// token, err
	_, err := s.authService.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			fmt.Fprintln(w, "Invalid username/password")
		default:
			fmt.Fprintln(w, "oops something went wrong")
		}

		return
	}

	fmt.Fprintln(w, "logged in")
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
