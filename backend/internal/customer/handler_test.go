package customer

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
	customers []Customer
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Customer, error) {
	out := make([]Customer, 0, len(f.customers))
	for _, c := range f.customers {
		if c.DeletedAt == nil {
			out = append(out, c)
		}
	}
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
	for i, cust := range f.customers {
		if cust.ID == c.ID {
			cust.LegalName = c.LegalName
			cust.TradeName = c.TradeName
			cust.DocumentID = c.DocumentID
			cust.Email = c.Email
			cust.Phone = c.Phone
			cust.UpdatedAt = c.UpdatedAt
			f.customers[i] = cust
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, c := range f.customers {
		if c.ID == id {
			if c.DeletedAt != nil {
				return sql.ErrNoRows
			}
			now := time.Now()
			c.DeletedAt = &now
			f.customers[i] = c
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

func TestUpdateCustomer(t *testing.T) {
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

	var created Customer
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	updBody := strings.NewReader(`{"legal_name":"ACME Updated","document_id":"1"}`)
	req, err := http.NewRequest(http.MethodPut, server.URL+"/customers/"+created.ID, updBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /customers/{id} error: %v", err)
	}
	defer resp2.Body.Close()

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

	if len(out) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(out))
	}

	if out[0].LegalName != "ACME Updated" {
		t.Fatalf("legal_name not updated: %+v", out[0])
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var created Customer
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/customers/"+created.ID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /customers/{id} error: %v", err)
	}
	defer resp2.Body.Close()

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
