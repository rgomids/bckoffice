package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type fakeRepository struct {
	services []Service
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Service, error) {
	out := make([]Service, 0, len(f.services))
	for _, s := range f.services {
		if s.DeletedAt == nil {
			out = append(out, s)
		}
	}
	return out, nil
}

func (f *fakeRepository) Create(ctx context.Context, s *Service) error {
	f.services = append(f.services, *s)
	return nil
}

func (f *fakeRepository) Update(ctx context.Context, s *Service) error {
	for i, svc := range f.services {
		if svc.ID == s.ID {
			svc.Name = s.Name
			svc.Description = s.Description
			svc.BasePrice = s.BasePrice
			svc.IsActive = s.IsActive
			svc.UpdatedAt = s.UpdatedAt
			f.services[i] = svc
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, svc := range f.services {
		if svc.ID == id {
			if svc.DeletedAt != nil {
				return sql.ErrNoRows
			}
			now := time.Now()
			svc.DeletedAt = &now
			f.services[i] = svc
			return nil
		}
	}
	return sql.ErrNoRows
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

func TestUpdateService(t *testing.T) {
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

	var created Service
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	upd := strings.NewReader(`{"name":"Site Avancado","base_price":2000}`)
	req, err := http.NewRequest(http.MethodPut, server.URL+"/services/"+created.ID, upd)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /services/{id} error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/services")
	if err != nil {
		t.Fatalf("GET /services error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Service
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 service, got %d", len(list))
	}
	if list[0].Name != "Site Avancado" {
		t.Fatalf("name not updated: %+v", list[0])
	}
}

func TestDeleteService(t *testing.T) {
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

	var created Service
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/services/"+created.ID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /services/{id} error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/services")
	if err != nil {
		t.Fatalf("GET /services error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Service
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("expected 0 services, got %d", len(list))
	}
}
