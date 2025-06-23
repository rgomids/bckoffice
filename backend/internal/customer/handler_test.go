package customer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"database/sql"
	"github.com/go-chi/chi/v5"
)

type fakeRepository struct {
	customers []Customer
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Customer, error) {
	out := make([]Customer, len(f.customers))
	copy(out, f.customers)
	return out, nil
}

func (f *fakeRepository) FindByID(ctx context.Context, id string) (Customer, error) {
	return Customer{}, nil
}

func (f *fakeRepository) Create(ctx context.Context, c *Customer, addresses []Address) error {
	f.customers = append(f.customers, *c)
	return nil
}

func (f *fakeRepository) Update(ctx context.Context, c *Customer, addresses []Address) error {
	return nil
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, c := range f.customers {
		if c.ID == id {
			f.customers = append(f.customers[:i], f.customers[i+1:]...)
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

func TestGetCustomersEmpty(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/customers")
	if err != nil {
		t.Fatalf("GET /customers error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out []Customer
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(out) != 0 {
		t.Fatalf("expected 0 customers, got %d", len(out))
	}
}

func TestPostCustomersAndGet(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"legal_name":"ACME","document_id":"1","addresses":[{"address_type":"billing","street":"A","city":"X","state":"Y"}]}`)
	resp, err := http.Post(server.URL+"/customers", "application/json", body)
	if err != nil {
		t.Fatalf("POST /customers error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	resp2, err := http.Get(server.URL + "/customers")
	if err != nil {
		t.Fatalf("GET /customers error: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}

	var out []Customer
	if err := json.NewDecoder(resp2.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(out))
	}

	if out[0].LegalName != "ACME" || out[0].DocumentID != "1" {
		t.Fatalf("unexpected customer data: %+v", out[0])
	}
}

func TestDeleteCustomer(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"legal_name":"ACME","document_id":"1","addresses":[{"address_type":"billing","street":"A","city":"X","state":"Y"}]}`)
	resp, err := http.Post(server.URL+"/customers", "application/json", body)
	if err != nil {
		t.Fatalf("POST /customers error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	id := strings.TrimPrefix(resp.Header.Get("Location"), "/customers/")
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/customers/"+id, nil)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /customers error: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/customers")
	if err != nil {
		t.Fatalf("GET /customers error: %v", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp3.StatusCode)
	}
	var out []Customer
	if err := json.NewDecoder(resp3.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 customers, got %d", len(out))
	}
}
