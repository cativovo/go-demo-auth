package main

import (
	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/cativovo/go-demo-auth/pkg/http"
	"github.com/cativovo/go-demo-auth/pkg/storage/postgres"
	"github.com/cativovo/go-demo-auth/pkg/storage/supabase"
	"github.com/cativovo/go-demo-auth/pkg/user"
)

func main() {
	supabaseRepository := supabase.NewSupabaseRepository()
	pgRepository := postgres.NewPostgresRepository()
	r := struct {
		*supabase.SupabaseRepository
		*postgres.PostgresRepository
	}{
		supabaseRepository,
		pgRepository,
	}

	authService := auth.NewAuthService(supabaseRepository)
	userService := user.NewUserService(r)

	server := http.NewServer(authService, userService)

	server.ListenAndServe("127.0.0.1:3000")
}
