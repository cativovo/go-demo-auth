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
	ErrInvalidEmail     = errors.New("invalid email")
	ErrUserNotFound     = errors.New("user not found")
)

type User struct {
	Id    string
	Email string
	Name  string
}

type Credentials struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required,min=6"`
}

type Service interface {
	Register(u Credentials) (auth.Token, []error)
	ValidateCredentials(c Credentials) validator.ValidationErrors
	GetUserByEmail(email string) (User, error)
}

type Repository interface {
	AddUser(u User) (User, error)
	Register(email, password string) (auth.Token, error)
	GetUserByEmail(email string) (User, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
}

func NewUserService(r Repository) Service {
	v := validator.New(validator.WithRequiredStructEnabled())

	return &service{
		repository: r,
		validate:   v,
	}
}

func (s *service) Register(c Credentials) (auth.Token, []error) {
	var errors []error

	if err := s.validate.Struct(c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errors = append(errors, fmt.Errorf("%s is %w", err.Field(), ErrFieldRequired))
			}

			switch err.StructField() {
			case "Password":
				errors = append(errors, ErrPasswordTooShort)
			case "Email":
				errors = append(errors, ErrInvalidEmail)
			}
		}
	}

	if errors != nil {
		return auth.Token{}, errors
	}

	token, err := s.repository.Register(c.Email, c.Password)
	if err != nil {
		log.Println("UserService Register repository.Register:", err)
		return auth.Token{}, append(errors, err)
	}

	user := User{
		Id:    token.UserId,
		Email: c.Email,
		Name:  c.Name,
	}

	if _, err := s.repository.AddUser(user); err != nil {
		log.Println("UserService Register AddUser:", err)
		return auth.Token{}, append(errors, err)
	}

	return token, nil
}

func (s *service) ValidateCredentials(c Credentials) validator.ValidationErrors {
	err := s.validate.Struct(c)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (s *service) GetUserByEmail(email string) (User, error) {
	u, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return User{}, ErrUserNotFound
	}

	return u, nil
}
