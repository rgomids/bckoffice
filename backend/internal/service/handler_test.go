package service

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
	services []Service
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Service, error) {
	out := make([]Service, len(f.services))
	copy(out, f.services)
	return out, nil
}

func (f *fakeRepository) Create(ctx context.Context, s *Service) error {
	f.services = append(f.services, *s)
	return nil
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	repo := &fakeRepository{}
	RegisterRoutes(r, repo)
	return r
}

func TestGetServicesEmpty(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/services")
	if err != nil {
		t.Fatalf("GET /services error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out []Service
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(out) != 0 {
		t.Fatalf("expected 0 services, got %d", len(out))
	}
}

func TestCreateService(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"name":"Site Basico","base_price":1500}`)
	resp, err := http.Post(server.URL+"/services", "application/json", body)
	if err != nil {
		t.Fatalf("POST /services error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var svc Service
	if err := json.NewDecoder(resp.Body).Decode(&svc); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if svc.Name != "Site Basico" || svc.BasePrice != 1500 {
		t.Fatalf("unexpected service data: %+v", svc)
	}
}
