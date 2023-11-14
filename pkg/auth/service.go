package auth

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

type Token struct {
	UserId       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	ExpiresAt    int
}

type Service interface {
	Login(email, password string) (Token, error)
	Logout(token string) error
}

type Repository interface {
	Login(email, password string) (Token, error)
	Logout(token string) error
}

type service struct {
	repository Repository
}

func NewAuthService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Login(email, password string) (Token, error) {
	return s.repository.Login(email, password)
}

func (s *service) Logout(token string) error {
	return s.repository.Logout(token)
}
