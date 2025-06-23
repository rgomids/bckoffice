package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"context"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/smithl4b/rcm.backoffice/internal/users"
)

type fakeRepository struct {
	user users.User
	err  error
}

func (f *fakeRepository) FindByEmail(ctx context.Context, email string) (users.User, error) {
	if f.err != nil {
		return users.User{}, f.err
	}
	if email == f.user.Email {
		return f.user, nil
	}
	return users.User{}, sql.ErrNoRows
}

func setupRouter(repo Repository) *chi.Mux {
	r := chi.NewRouter()
	RegisterRoutes(r, repo)
	r.Group(func(pr chi.Router) {
		pr.Use(AuthMiddleware)
		pr.Get("/private", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	})
	return r
}

func TestLoginSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo := &fakeRepository{user: users.User{ID: "1", Email: "foo@example.com", PasswordHash: string(hash), Role: "admin"}}

	server := httptest.NewServer(setupRouter(repo))
	defer server.Close()

	body := strings.NewReader(`{"email":"foo@example.com","password":"pass"}`)
	resp, err := http.Post(server.URL+"/login", "application/json", body)
	if err != nil {
		t.Fatalf("POST /login error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if out.Token == "" {
		t.Fatalf("expected token, got empty")
	}
}

func TestLoginFailure(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo := &fakeRepository{user: users.User{ID: "1", Email: "foo@example.com", PasswordHash: string(hash), Role: "admin"}}
	server := httptest.NewServer(setupRouter(repo))
	defer server.Close()

	body := strings.NewReader(`{"email":"foo@example.com","password":"wrong"}`)
	resp, err := http.Post(server.URL+"/login", "application/json", body)
	if err != nil {
		t.Fatalf("POST /login error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo := &fakeRepository{user: users.User{ID: "1", Email: "foo@example.com", PasswordHash: string(hash), Role: "admin"}}
	server := httptest.NewServer(setupRouter(repo))
	defer server.Close()

	body := strings.NewReader(`{"email":"foo@example.com","password":"pass"}`)
	resp, err := http.Post(server.URL+"/login", "application/json", body)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	defer resp.Body.Close()
	var out AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode login: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/private", nil)
	req.Header.Set("Authorization", "Bearer "+out.Token)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /private: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}

	req2, _ := http.NewRequest(http.MethodGet, server.URL+"/private", nil)
	resp3, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("GET /private: %v", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp3.StatusCode)
	}
}
