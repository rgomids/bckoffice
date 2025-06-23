package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeRepo struct {
	logs []AuditLog
}

func (f *fakeRepo) Create(_ context.Context, l *AuditLog) error {
	f.logs = append(f.logs, *l)
	return nil
}

type fakeGeo struct{}

func (f fakeGeo) Lookup(_ context.Context, ip string) (GeoInfo, error) {
	return GeoInfo{Country: "BR", City: "Sao Paulo", Lat: 0, Lon: 0}, nil
}

func TestMiddlewareInsert(t *testing.T) {
	repo := &fakeRepo{}
	geo := fakeGeo{}
	r := chi.NewRouter()
	r.Use(NewAuditMiddleware(repo, geo))
	r.Post("/customers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	body := strings.NewReader(`{"name":"x"}`)
	resp, err := http.Post(server.URL+"/customers", "application/json", body)
	if err != nil {
		t.Fatalf("POST request error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
	if len(repo.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(repo.logs))
	}
	if repo.logs[0].Action != "insert" {
		t.Fatalf("unexpected action: %s", repo.logs[0].Action)
	}
}
