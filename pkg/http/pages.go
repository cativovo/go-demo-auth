package http

import (
	"html/template"
	"net/http"
)

var registerPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/components/register_form.html",
	),
)

var loginPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/components/login_form.html",
	),
)

var accountPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/components/account.html",
	),
)

func (s *Server) registerPages() {
	s.router.Get("/", s.loginPage)
	s.router.Get("/register", s.registerPage)
	s.router.Get("/account", s.accountPage)
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	loginPageTmpl.Execute(w, nil)
}

func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	registerPageTmpl.Execute(w, nil)
}

func (s *Server) accountPage(w http.ResponseWriter, r *http.Request) {
	accountPageTmpl.Execute(w, nil)
}
