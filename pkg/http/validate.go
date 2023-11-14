package http

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/cativovo/go-demo-auth/pkg/user"
)

var (
	registerFormTmpl *template.Template = template.Must(template.ParseFiles("web/components/register_form.html"))
	errorAlertTmpl   *template.Template = template.Must(template.ParseFiles("web/components/error_alert.html"))
)

func (s *Server) registerValidateRoutes() {
	s.router.Post("/validate-register", s.handleValidateRegister)
}

func (s *Server) handleValidateRegister(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	name := strings.TrimSpace(r.FormValue("name"))

	data := map[string]any{
		"Email":    email,
		"Password": password,
		"Name":     name,
	}

	errs := s.userService.ValidateCredentials(user.Credentials{
		Email:    email,
		Password: password,
		Name:     name,
	})

	for _, err := range errs {
		switch err.StructField() {
		case "Email":
			data["ErrEmail"] = "Invalid Email!"
		case "Password":
			data["ErrPassword"] = "Password Too Short!"
		}
	}

	if data["ErrEmail"] != nil {
		registerFormTmpl.Execute(w, data)
		return
	}

	_, err := s.userService.GetUserByEmail(email)

	if err == nil {
		data["ErrEmail"] = "Email is already used!"
		registerFormTmpl.Execute(w, data)
		return
	}

	if !errors.Is(err, user.ErrUserNotFound) {
		errorAlertTmpl.Execute(w, map[string]any{
			"Message": "Something went wrong...",
		})
		return
	}

	data["AreValuesValid"] = errs == nil

	registerFormTmpl.Execute(w, data)
}
