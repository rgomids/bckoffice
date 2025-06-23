package finance

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smithl4b/rcm.backoffice/internal/auth"
)

type fakeRepository struct {
	receivables []AccountReceivable
	commissions []Commission
}

func (f *fakeRepository) ListReceivables(ctx context.Context, status string) ([]AccountReceivable, error) {
	out := make([]AccountReceivable, 0)
	for _, ar := range f.receivables {
		if status == "" || ar.Status == status {
			out = append(out, ar)
		}
	}
	return out, nil
}

func (f *fakeRepository) MarkAsPaid(ctx context.Context, id string) error {
	for i, ar := range f.receivables {
		if ar.ID == id {
			if ar.Status != "open" {
				return ErrAlreadyPaid
			}
			ar.Status = "paid"
			now := time.Now()
			ar.PaidAt = &now
			f.receivables[i] = ar
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) ListCommissions(ctx context.Context, onlyPending bool) ([]Commission, error) {
	out := make([]Commission, 0)
	for _, c := range f.commissions {
		if !onlyPending || !c.Approved {
			out = append(out, c)
		}
	}
	return out, nil
}

func (f *fakeRepository) ApproveCommission(ctx context.Context, id string, approverID string) error {
	for i, c := range f.commissions {
		if c.ID == id {
			if c.Approved {
				return ErrAlreadyApproved
			}
			c.Approved = true
			c.ApprovedBy = approverID
			now := time.Now()
			c.ApprovedAt = &now
			f.commissions[i] = c
			return nil
		}
	}
	return sql.ErrNoRows
}

func setupRouter(repo Repository) (*chi.Mux, string) {
	os.Setenv("JWT_SECRET", "testsecret")
	claims := jwt.MapClaims{"sub": "u1", "role": "finance", "exp": time.Now().Add(time.Hour).Unix()}
	tokenStr, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("testsecret"))

	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware)
	RegisterRoutes(r, repo)
	return r, tokenStr
}

func TestListReceivablesEmpty(t *testing.T) {
	router, token := setupRouter(&fakeRepository{})
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/receivables", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /receivables error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out []AccountReceivable
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}

func TestMarkAsPaid(t *testing.T) {
	repo := &fakeRepository{receivables: []AccountReceivable{{ID: "r1", Status: "open"}}}
	router, token := setupRouter(repo)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/receivables/r1/pay", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /receivables/{id}/pay error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.StatusCode)
	}

	if repo.receivables[0].Status != "paid" {
		t.Fatalf("status not updated: %+v", repo.receivables[0])
	}
	if repo.receivables[0].PaidAt == nil {
		t.Fatalf("paid_at not set")
	}
}

func TestApproveCommission(t *testing.T) {
	repo := &fakeRepository{commissions: []Commission{{ID: "c1"}}}
	router, token := setupRouter(repo)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/commissions/c1/approve", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /commissions/{id}/approve error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.StatusCode)
	}

	if !repo.commissions[0].Approved {
		t.Fatalf("commission not approved")
	}
	if repo.commissions[0].ApprovedBy == "" {
		t.Fatalf("approved_by empty")
	}
}
