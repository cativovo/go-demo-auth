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

type (
	Json    map[string]interface{}
	Headers map[string]string
)

type Response struct {
	data Json
	*http.Response
}

func NewAuthService(project, apiKey string) *AuthService {
	baseUrl := fmt.Sprintf("https://%s.supabase.co/auth/v1", project)
	return &AuthService{
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}

func (a *AuthService) fetch(method string, path string, headers Headers, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(method, a.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apiKey", a.apiKey)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

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

	data := Json{}

	if err := json.Unmarshal([]byte(r), &data); err != nil {
		return nil, err
	}

	return &Response{
			data:     data,
			Response: res,
		},
		nil
}

func (a *AuthService) Login(userCredentials UserCredentials) (*Response, error) {
	payload, err := json.Marshal(userCredentials)
	if err != nil {
		return nil, err
	}

	res, err := a.fetch("POST", "/token?grant_type=password", nil, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *AuthService) Logout(token string) error {
	headers := Headers{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	_, err := a.fetch("POST", "/logout?scope=local", headers, nil)
	if err != nil {
		return nil
	}

	return nil
}

func (a *AuthService) Register(userCredentials UserCredentials) (*Response, error) {
	payload, err := json.Marshal(userCredentials)
	if err != nil {
		return nil, err
	}

	res, err := a.fetch("POST", "/signup", nil, bytes.NewBuffer(payload))
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

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		userCredentials := UserCredentials{
			Email:    r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}

		_, err := authService.Register(userCredentials)
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, "request error")
			return
		}

		// if res.StatusCode == http.StatusBadRequest {
		// 	fmt.Fprintln(w, "Invalid email/password")
		// 	return
		// }

		fmt.Fprintln(w, "registered")
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		userCredentials := UserCredentials{
			Email:    r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}

		res, err := authService.Login(userCredentials)
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, "request error")
			return
		}

		if res.StatusCode == http.StatusBadRequest {
			fmt.Fprintln(w, "Invalid email/password")
			return
		}

		if v, ok := res.data["access_token"].(string); ok {
			cookie := &http.Cookie{
				Name:  "jwt",
				Value: v,
			}

			http.SetCookie(w, cookie)
		}

		fmt.Fprintln(w, "logged in")
	})

	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		jwt, err := r.Cookie("jwt")
		if err != nil {
			return
		}

		if err := authService.Logout(jwt.Value); err != nil {
			fmt.Fprintln(w, "ooo")
			return
		}
	})

	if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}
