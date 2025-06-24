package auditquery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rgomids/bckoffice/internal/audit"
	"github.com/rgomids/bckoffice/internal/auth"
)

type fakeRepo struct {
	logs []audit.AuditLog
}

func (f *fakeRepo) List(ctx context.Context, fl AuditFilter) ([]audit.AuditLog, error) {
	out := make([]audit.AuditLog, 0)
	for _, l := range f.logs {
		if fl.EntityName != "" && l.EntityName != fl.EntityName {
			continue
		}
		if fl.Action != "" && l.Action != fl.Action {
			continue
		}
		out = append(out, l)
	}
	if fl.Limit > 0 && len(out) > fl.Limit {
		out = out[:fl.Limit]
	}
	return out, nil
}

func setupRouter(repo Repository, role string) (*chi.Mux, string) {
	os.Setenv("JWT_SECRET", "testsecret")
	claims := jwt.MapClaims{"sub": "u1", "role": role, "exp": time.Now().Add(time.Hour).Unix()}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("testsecret"))

	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware)
	RegisterRoutes(r, repo)
	return r, token
}

func TestListAuditLogsFilter(t *testing.T) {
	repo := &fakeRepo{logs: []audit.AuditLog{
		{ID: "1", EntityName: "customers", Action: "update"},
		{ID: "2", EntityName: "services", Action: "insert"},
	}}
	r, token := setupRouter(repo, "admin")
	server := httptest.NewServer(r)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/audit-logs?entity=customers&action=update&limit=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var out []audit.AuditLog
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 1 || out[0].ID != "1" {
		t.Fatalf("unexpected result: %+v", out)
	}
}
