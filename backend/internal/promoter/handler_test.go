package promoter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeRepository struct {
	promoters []Promoter
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Promoter, error) {
	out := make([]Promoter, 0, len(f.promoters))
	for _, p := range f.promoters {
		if p.DeletedAt == nil {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *fakeRepository) Create(ctx context.Context, p *Promoter) error {
	f.promoters = append(f.promoters, *p)
	return nil
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	repo := &fakeRepository{}
	RegisterRoutes(r, repo)
	return r
}

func TestGetPromotersEmpty(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/promoters")
	if err != nil {
		t.Fatalf("GET /promoters error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out []Promoter
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(out) != 0 {
		t.Fatalf("expected 0 promoters, got %d", len(out))
	}
}

func TestCreatePromoter(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"full_name":"Joao","email":"joao@ex.com"}`)
	resp, err := http.Post(server.URL+"/promoters", "application/json", body)
	if err != nil {
		t.Fatalf("POST /promoters error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var p Promoter
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if p.FullName != "Joao" || p.Email != "joao@ex.com" {
		t.Fatalf("unexpected promoter data: %+v", p)
	}
}
