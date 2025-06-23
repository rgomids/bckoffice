package promoter

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

func (f *fakeRepository) Update(ctx context.Context, p *Promoter) error {
	for i, pr := range f.promoters {
		if pr.ID == p.ID {
			pr.FullName = p.FullName
			pr.Email = p.Email
			pr.Phone = p.Phone
			pr.DocumentID = p.DocumentID
			pr.BankAccount = p.BankAccount
			pr.UpdatedAt = p.UpdatedAt
			f.promoters[i] = pr
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, pr := range f.promoters {
		if pr.ID == id {
			if pr.DeletedAt != nil {
				return sql.ErrNoRows
			}
			now := time.Now()
			pr.DeletedAt = &now
			f.promoters[i] = pr
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

func TestUpdatePromoter(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"full_name":"Joao"}`)
	resp, err := http.Post(server.URL+"/promoters", "application/json", body)
	if err != nil {
		t.Fatalf("POST /promoters error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var created Promoter
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	upd := strings.NewReader(`{"full_name":"Maria"}`)
	req, err := http.NewRequest(http.MethodPut, server.URL+"/promoters/"+created.ID, upd)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /promoters/{id} error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/promoters")
	if err != nil {
		t.Fatalf("GET /promoters error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Promoter
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 promoter, got %d", len(list))
	}
	if list[0].FullName != "Maria" {
		t.Fatalf("name not updated: %+v", list[0])
	}
}

func TestDeletePromoter(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"full_name":"Joao"}`)
	resp, err := http.Post(server.URL+"/promoters", "application/json", body)
	if err != nil {
		t.Fatalf("POST /promoters error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var created Promoter
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/promoters/"+created.ID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /promoters/{id} error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/promoters")
	if err != nil {
		t.Fatalf("GET /promoters error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Promoter
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("expected 0 promoters, got %d", len(list))
	}
}
