package postgres

import (
	"context"
	"log"
	"os"

	postgres "github.com/cativovo/go-demo-auth/pkg/storage/postgres/sqlc_generated"
	"github.com/cativovo/go-demo-auth/pkg/user"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	ctx     context.Context
	queries *postgres.Queries
}

func NewPostgresRepository() *PostgresRepository {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}

	queries := postgres.New(conn)

	return &PostgresRepository{
		ctx:     ctx,
		queries: queries,
	}
}

func (r *PostgresRepository) AddUser(u user.User) (user.User, error) {
	p := postgres.AddUserParams{
		ID:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}

	newUser, err := r.queries.AddUser(r.ctx, p)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		Id:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}

func (r *PostgresRepository) GetUserByEmail(email string) (user.User, error) {
	u, err := r.queries.GetUserByEmail(r.ctx, email)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		Id:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}
