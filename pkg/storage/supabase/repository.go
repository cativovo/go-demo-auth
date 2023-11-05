package supabase

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cativovo/go-demo-auth/pkg/auth"
	"github.com/cativovo/go-demo-auth/pkg/user"
)

type SupabaseRepository struct {
	httpClient *http.Client
	apiKey     string
	baseUrl    string
}

type headers map[string]string

type response struct {
	*http.Response
	Data map[string]any
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

	res, err := s.fetch("POST", "/signup", nil, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Supabase Register:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	if res.StatusCode == http.StatusBadRequest {
		return auth.Token{}, user.ErrEmailAlreadyUsed
	}

	if res.StatusCode > http.StatusBadRequest {
		log.Printf("Supabase Register status code (%d)", res.StatusCode)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	token, err := s.getToken(res)
	if err != nil {
		log.Println("Supabase Register getToken:", err)
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	return token, nil
}

func (s *SupabaseRepository) Login(email, password string) (auth.Token, error) {
	c := credentials{
		Email:    email,
		Password: password,
	}
	payload, err := json.Marshal(c)
	if err != nil {
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	res, err := s.fetch("POST", "/token?grant_type=password", nil, bytes.NewBuffer(payload))
	if err != nil {
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	if res.StatusCode == http.StatusBadRequest {
		return auth.Token{}, auth.ErrInvalidCredentials
	}

	if res.StatusCode > http.StatusBadRequest {
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	token, err := s.getToken(res)
	if err != nil {
		return auth.Token{}, auth.ErrSomethingWentWrong
	}

	return token, nil
}

func (s *SupabaseRepository) Logout(token string) error {
	headers := headers{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	res, err := s.fetch("POST", "/logout?scope=local", headers, nil)
	if err != nil {
		return auth.ErrSomethingWentWrong
	}

	if res.StatusCode != http.StatusNoContent {
		return auth.ErrSomethingWentWrong
	}

	return nil
}

// helpers
func (s *SupabaseRepository) fetch(method string, path string, h headers, body io.Reader) (*response, error) {
	req, err := http.NewRequest(method, s.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apiKey", s.apiKey)

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var r string

	for scanner.Scan() {
		r += scanner.Text()
	}

	data := map[string]any{}

	if err := json.Unmarshal([]byte(r), &data); err != nil {
		return nil, err
	}

	return &response{
			Data:     data,
			Response: res,
		},
		nil
}

func (s *SupabaseRepository) getToken(res *response) (auth.Token, error) {
	accessToken, ok := res.Data["access_token"].(string)
	if !ok {
		return auth.Token{}, errors.New("access_token field not found")
	}

	refreshToken, ok := res.Data["refresh_token"].(string)
	if !ok {
		return auth.Token{}, errors.New("refresh_token field not found")
	}

	expiresIn, ok := res.Data["expires_in"].(float64)
	if !ok {
		return auth.Token{}, errors.New("expires_in field not found")
	}

	expiresAt, ok := res.Data["expires_at"].(float64)
	if !ok {
		return auth.Token{}, errors.New("expires_at field not found")
	}

	var userId string
	if user, ok := res.Data["user"].(map[string]any); ok {
		if id, ok := user["id"].(string); ok {
			userId = id
		} else {
			return auth.Token{}, auth.ErrSomethingWentWrong
		}
	}

	return auth.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
			ExpiresIn:    expiresIn,
			UserId:       userId,
		},
		nil
}
