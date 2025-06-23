package contract

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
