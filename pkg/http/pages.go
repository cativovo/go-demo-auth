package http

import (
	"html/template"
	"net/http"
)

var baseTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/components/register_form.html",
		"web/components/email_input.html",
	),
)

func (s *Server) registerPages() {
	s.router.Get("/", s.homePage)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	baseTmpl.Execute(w, nil)
}
