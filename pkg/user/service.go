package user

import (
	"errors"
	"fmt"
	"log"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/go-playground/validator/v10"
)

var (
	ErrEmailAlreadyUsed = errors.New("email is already used")
	ErrFieldRequired    = errors.New("required")
	ErrPasswordTooShort = errors.New("password is too short")
)

type User struct {
	Id    string
	Email string
	Name  string
}

type Credentials struct {
	Email    string `validate:"required,email" json:"Foo"`
	Name     string `validate:"required"`
	Password string `validate:"required,min=6"`
}

type Service interface {
	Register(u Credentials) (auth.Token, error)
}

type Repository interface {
	AddUser(u User) (User, error)
	Register(email, password string) (auth.Token, error)
}

type service struct {
	r Repository
}

var validate *validator.Validate

func NewUserService(r Repository) Service {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return &service{
		r: r,
	}
}

func (s *service) Register(u Credentials) (auth.Token, error) {
	if err := validate.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return auth.Token{}, fmt.Errorf("%s is %w", err.Field(), ErrFieldRequired)
			}

			switch err.StructField() {
			case "Password":
				return auth.Token{}, ErrPasswordTooShort
			}
		}

		return auth.Token{}, err
	}

	token, err := s.r.Register(u.Email, u.Password)
	if err != nil {
		log.Println("UserService Register:", err)
		return auth.Token{}, err
	}

	user := User{
		Id:    token.UserId,
		Email: u.Email,
		Name:  u.Name,
	}
	s.r.AddUser(user)

	return auth.Token{}, nil
}
