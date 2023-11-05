package postgres

import (
	"log"

	"github.com/cativovo/go-demo-auth/pkg/user"
)

type PostgresRepository struct{}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) AddUser(u user.User) (user.User, error) {
	log.Println("User added!", u)
	return u, nil
}
