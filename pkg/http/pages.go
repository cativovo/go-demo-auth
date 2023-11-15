package http

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var registerPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/layouts/public.html",
		"web/components/register_form.html",
	),
)

var loginPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/layouts/public.html",
		"web/components/login_form.html",
	),
)

var accountPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/layouts/private.html",
		"web/components/nav.html",
		"web/components/account.html",
	),
)

var infoPageTmpl *template.Template = template.Must(
	template.ParseFiles(
		"web/base.html",
		"web/layouts/private.html",
		"web/components/nav.html",
		"web/components/info.html",
	),
)

func (s *Server) registerPages() {
	s.router.Route("/auth-page", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := r.Cookie("access_token")
				if err == nil {
					w.Header().Add("Location", "/")
					w.WriteHeader(http.StatusFound)
					return
				}

				next.ServeHTTP(w, r)
			})
		})
		r.Get("/login", s.loginPage)
		r.Get("/register", s.registerPage)
	})

	s.router.Route("/", func(r chi.Router) {
		r.Use(authMiddleWare(s.authService))
		r.Get("/", s.accountPage)
		r.Get("/info", s.infoPage)
	})
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-store, public")
	loginPageTmpl.Execute(w, nil)
}

func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-store, public")
	registerPageTmpl.Execute(w, nil)
}

func (s *Server) accountPage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdKey).(string)

	user, err := s.userService.GetUserById(userId)
	if err != nil {
		log.Println(err)
		w.Header().Add("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	data := map[string]any{
		"UserId": user.Id,
		"Name":   user.Name,
	}

	w.Header().Add("Cache-Control", "no-store, private")

	if r.Header.Get("HX-Boosted") == "true" {
		accountPageTmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	accountPageTmpl.Execute(w, data)
}

func (s *Server) infoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "private, max-age=30")

	if r.Header.Get("HX-Boosted") == "true" {
		infoPageTmpl.ExecuteTemplate(w, "layout", nil)
		return
	}

	infoPageTmpl.Execute(w, nil)
}
