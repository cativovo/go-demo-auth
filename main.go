package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var httpClient *http.Client

type AuthService struct {
	apiKey  string
	baseUrl string
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JsonResponse map[string]interface{}

func NewAuthService(project, apiKey string) *AuthService {
	baseUrl := fmt.Sprintf("https://%s.supabase.co/auth/v1", project)
	return &AuthService{
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}

func (a *AuthService) fetch(method string, path string, body io.Reader) (JsonResponse, error) {
	req, err := http.NewRequest(method, a.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apiKey", a.apiKey)

	res, err := httpClient.Do(req)
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

	j := JsonResponse{}

	if err := json.Unmarshal([]byte(r), &j); err != nil {
		return nil, err
	}

	return j, nil
}

func (a *AuthService) Login(email, password string) (JsonResponse, error) {
	userCredentials := UserCredentials{
		Email:    email,
		Password: password,
	}
	payload, err := json.Marshal(userCredentials)
	if err != nil {
		return nil, err
	}

	res, err := a.fetch("POST", "/token?grant_type=password", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	r := chi.NewRouter()
	httpClient = &http.Client{}

	authService := NewAuthService(os.Getenv("SUPABASE_PROJECT"), os.Getenv("SUPABASE_API_KEY"))

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := authService.Login("test@example.com", "1234")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, "request error")
		}

		j, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, "json Marshal error")
		}

		log.Println(res["access_token"])
		fmt.Fprintln(w, string(j))
	})

	http.ListenAndServe("127.0.0.1:3000", r)
}
