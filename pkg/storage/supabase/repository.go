package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	userService "github.com/cativovo/go-demo-auth/pkg/user"
)

type SupabaseRepository struct {
	httpClient *http.Client
	apiKey     string
	baseUrl    string
}

type user struct {
	Id string `json:"id"`
}

type token struct {
	User         user   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSupabaseRepository() *SupabaseRepository {
	return &SupabaseRepository{
		apiKey:     os.Getenv("SUPABASE_API_KEY"),
		baseUrl:    fmt.Sprintf("https://%s.supabase.co/auth/v1", os.Getenv("SUPABASE_PROJECT")),
		httpClient: &http.Client{},
	}
}

func (s *SupabaseRepository) Register(email, password string) (auth.Token, error) {
	c := credentials{
		Email:    email,
		Password: password,
	}

	payload, err := json.Marshal(c)
	if err != nil {
		log.Println("Supabase Register:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	req, err := s.newRequest("POST", "/signup", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Supabase Register newRequest:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Println("Supabase Register Do:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	if res.StatusCode == http.StatusBadRequest {
		return auth.Token{}, userService.ErrEmailAlreadyUsed
	}

	t := token{}

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Println("Supabase Register Decode:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	return auth.Token{
			UserId:       t.User.Id,
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
			ExpiresIn:    t.ExpiresIn,
			ExpiresAt:    t.ExpiresAt,
		},
		nil
}

func (s *SupabaseRepository) Login(email, password string) (auth.Token, error) {
	c := credentials{
		Email:    email,
		Password: password,
	}

	payload, err := json.Marshal(c)
	if err != nil {
		log.Println("Supabase Login:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	req, err := s.newRequest("POST", "/token?grant_type=password", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Supabase Login newRequest:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Println("Supabase Login Do:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	if res.StatusCode == http.StatusBadRequest {
		return auth.Token{}, auth.ErrInvalidCredentials
	}

	t := token{}

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Println("Supabase Login Decode:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	return auth.Token{
			UserId:       t.User.Id,
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
			ExpiresIn:    t.ExpiresIn,
			ExpiresAt:    t.ExpiresAt,
		},
		nil
}

func (s *SupabaseRepository) Logout(token string) error {
	req, err := s.newRequest("POST", "/logout?scope=local", nil)
	if err != nil {
		log.Println("Supabase Logout newRequest", err)
		return auth.ErrSomethingWentWrong
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Println("Supabase Logout Do", err)
		return auth.ErrSomethingWentWrong
	}

	if res.StatusCode != http.StatusNoContent {
		return auth.ErrSomethingWentWrong
	}

	return nil
}

func (s *SupabaseRepository) GetUserId(token string) (string, error) {
	req, err := s.newRequest("GET", "/user", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		log.Println("Supabase GetUserId newRequest", err)
		return "", auth.ErrSomethingWentWrong
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Println("Supabase GetUserId Do", err)
		return "", auth.ErrSomethingWentWrong
	}

	u := user{}

	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		log.Println("Supabase GetUserId Do", err)
		return "", auth.ErrSomethingWentWrong
	}

	return u.Id, nil
}

// helpers
func (s *SupabaseRepository) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, s.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apiKey", s.apiKey)

	return req, nil
}
