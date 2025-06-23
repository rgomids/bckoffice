package contract

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
	contracts []Contract
}

func (f *fakeRepository) FindAll(ctx context.Context) ([]Contract, error) {
	out := make([]Contract, 0, len(f.contracts))
	for _, c := range f.contracts {
		if c.DeletedAt == nil {
			out = append(out, c)
		}
	}
	return out, nil
}

func (f *fakeRepository) Create(ctx context.Context, c *Contract) error {
	f.contracts = append(f.contracts, *c)
	return nil
}

func (f *fakeRepository) Update(ctx context.Context, c *Contract) error {
	for i, ct := range f.contracts {
		if ct.ID == c.ID {
			ct.ValueTotal = c.ValueTotal
			ct.StartDate = c.StartDate
			ct.EndDate = c.EndDate
			ct.Status = c.Status
			ct.UpdatedAt = c.UpdatedAt
			f.contracts[i] = ct
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, ct := range f.contracts {
		if ct.ID == id {
			if ct.DeletedAt != nil {
				return sql.ErrNoRows
			}
			now := time.Now()
			ct.DeletedAt = &now
			f.contracts[i] = ct
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

func TestGetContractsEmpty(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	resp, err := http.Get(server.URL + "/contracts")
	if err != nil {
		t.Fatalf("GET /contracts error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out []Contract
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(out) != 0 {
		t.Fatalf("expected 0 contracts, got %d", len(out))
	}
}

func TestCreateContract(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"customer_id":"c1","service_id":"s1","value_total":1000,"start_date":"2025-07-01"}`)
	resp, err := http.Post(server.URL+"/contracts", "application/json", body)
	if err != nil {
		t.Fatalf("POST /contracts error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var c Contract
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if c.CustomerID != "c1" || c.ServiceID != "s1" || c.ValueTotal != 1000 {
		t.Fatalf("unexpected contract data: %+v", c)
	}

	if c.StartDate.IsZero() {
		t.Fatalf("start_date not set")
	}
}

func TestUpdateContract(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"customer_id":"c1","service_id":"s1","value_total":1000,"start_date":"2025-07-01"}`)
	resp, err := http.Post(server.URL+"/contracts", "application/json", body)
	if err != nil {
		t.Fatalf("POST /contracts error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var created Contract
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	upd := strings.NewReader(`{"value_total":2000,"start_date":"2025-07-01","status":"active"}`)
	req, err := http.NewRequest(http.MethodPut, server.URL+"/contracts/"+created.ID, upd)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /contracts/{id} error: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/contracts")
	if err != nil {
		t.Fatalf("GET /contracts error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Contract
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 contract, got %d", len(list))
	}
	if list[0].ValueTotal != 2000 {
		t.Fatalf("value_total not updated: %+v", list[0])
	}
}

func TestDeleteContract(t *testing.T) {
	server := httptest.NewServer(setupRouter())
	defer server.Close()

	body := strings.NewReader(`{"customer_id":"c1","service_id":"s1","value_total":1000,"start_date":"2025-07-01"}`)
	resp, err := http.Post(server.URL+"/contracts", "application/json", body)
	if err != nil {
		t.Fatalf("POST /contracts error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
	var created Contract
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/contracts/"+created.ID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /contracts/{id} error: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp2.StatusCode)
	}

	resp3, err := http.Get(server.URL + "/contracts")
	if err != nil {
		t.Fatalf("GET /contracts error: %v", err)
	}
	defer resp3.Body.Close()

	var list []Contract
	if err := json.NewDecoder(resp3.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 contracts, got %d", len(list))
	}
}
