package http

import (
	"context"
	"log"
	"net/http"

	"github.com/cativovo/go-demo-auth/pkg/auth"
)

type UserIdKey string

var userIdKey UserIdKey = "userId"

func setHtmlContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		next.ServeHTTP(w, r)
	})
}

func authMiddleWare(a auth.Service) func(next http.Handler) http.Handler {
	loginUrl := "/auth-page/login"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				log.Println(err)
				w.Header().Add("Location", loginUrl)
				w.WriteHeader(http.StatusFound)
				return
			}

			userId, err := a.GetUserId(accessTokenCookie.Value)
			if err != nil {
				log.Println(err)
				w.Header().Add("Location", loginUrl)
				w.WriteHeader(http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), userIdKey, userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
