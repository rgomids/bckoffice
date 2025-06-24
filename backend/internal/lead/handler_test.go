package lead

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smithl4b/rcm.backoffice/internal/auth"
)

type fakeRepository struct {
	leads []Lead
}

func (f *fakeRepository) List(ctx context.Context, status string) ([]Lead, error) {
	out := make([]Lead, 0)
	for _, l := range f.leads {
		if l.DeletedAt == nil && (status == "" || l.Status == status) {
			out = append(out, l)
		}
	}
	return out, nil
}

func (f *fakeRepository) Create(ctx context.Context, l *Lead) error {
	f.leads = append(f.leads, *l)
	return nil
}

func (f *fakeRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	for i, l := range f.leads {
		if l.ID == id {
			if allowedTransitions[l.Status] != newStatus {
				return errors.New("invalid status transition")
			}
			l.Status = newStatus
			f.leads[i] = l
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) Update(ctx context.Context, l *Lead) error {
	for i, lead := range f.leads {
		if lead.ID == l.ID {
			lead.ServiceID = l.ServiceID
			lead.PromoterID = l.PromoterID
			lead.Notes = l.Notes
			f.leads[i] = lead
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	for i, l := range f.leads {
		if l.ID == id {
			if l.DeletedAt != nil {
				return sql.ErrNoRows
			}
			now := time.Now()
			l.DeletedAt = &now
			f.leads[i] = l
			return nil
		}
	}
	return sql.ErrNoRows
}

func setupRouter(repo Repository, role string) (*chi.Mux, string) {
	os.Setenv("JWT_SECRET", "testsecret")
	claims := jwt.MapClaims{"sub": "u1", "role": role, "exp": time.Now().Add(time.Hour).Unix()}
	tokenStr, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("testsecret"))

	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware)
	RegisterRoutes(r, repo)
	return r, tokenStr
}

func TestCreateLead(t *testing.T) {
	repo := &fakeRepository{}
	router := chi.NewRouter()
	RegisterRoutes(router, repo)
	server := httptest.NewServer(router)
	defer server.Close()

	body := strings.NewReader(`{"customer_id":"c1","service_id":"s1"}`)
	resp, err := http.Post(server.URL+"/leads", "application/json", body)
	if err != nil {
		t.Fatalf("POST /leads error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var l Lead
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if l.CustomerID != "c1" || l.ServiceID != "s1" || l.Status != "lead" {
		t.Fatalf("unexpected lead data: %+v", l)
	}
}

func TestStatusTransitionInvalid(t *testing.T) {
	repo := &fakeRepository{leads: []Lead{{ID: "l1", Status: "lead"}}}
	r, token := setupRouter(repo, "finance")
	server := httptest.NewServer(r)
	defer server.Close()

	body := strings.NewReader(`{"status":"contract"}`)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/leads/l1/status", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /leads/{id}/status error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
