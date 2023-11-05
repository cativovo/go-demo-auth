package auth

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

type Token struct {
	UserId       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    float64
	ExpiresAt    float64
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
	r Repository
}

type headers map[string]string

type response struct {
	*http.Response
	Data map[string]any
}

func NewAuthService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) Login(email, password string) (Token, error) {
	return s.r.Login(email, password)
}

func (s *service) Logout(token string) error {
	return s.r.Logout(token)
}
