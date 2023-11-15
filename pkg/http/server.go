package http

import (
	"net/http"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/cativovo/go-demo-auth/pkg/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router      *chi.Mux
	authService auth.Service
	userService user.Service
}

func NewServer(a auth.Service, u user.Service) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(setHtmlContentTypeMiddleware)
	router.Use(middleware.Compress(5, "text/html", "text/css"))

	server := &Server{
		router:      router,
		authService: a,
		userService: u,
	}

	server.registerAuthRoutes()
	server.registerValidateRoutes()
	server.registerPages()

	return server
}

func (s *Server) ListenAndServe(addr string) {
	http.ListenAndServe(addr, s.router)
}
