package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"

	"github.com/rgomids/bckoffice/internal/auth"
)

// statusRecorder ajuda a capturar o status HTTP retornado.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// NewAuditMiddleware cria um middleware de auditoria.
func NewAuditMiddleware(repo Repository, geoSvc GeoService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}
			var bodyCopy []byte
			if r.Body != nil {
				bodyCopy, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewReader(bodyCopy))
			}
			sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sr, r)
			if sr.status < 200 || sr.status >= 300 {
				return
			}

			actionMap := map[string]string{
				http.MethodPost:   "insert",
				http.MethodPut:    "update",
				http.MethodDelete: "delete",
			}
			action := actionMap[r.Method]
			entity := r.Header.Get("X-Entity")
			if entity == "" {
				path := strings.TrimPrefix(r.URL.Path, "/")
				parts := strings.Split(path, "/")
				if len(parts) > 0 {
					entity = parts[0]
				}
			}
			entID := chi.URLParam(r, "id")
			userID := auth.UserIDFromContext(r.Context())
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = strings.Split(r.RemoteAddr, ":")[0]
			}
			ua := r.UserAgent()
			geoInfo, _ := geoSvc.Lookup(context.Background(), ip)
			geoBytes, _ := json.Marshal(geoInfo)
			var diff json.RawMessage
			if action == "update" && len(bodyCopy) > 0 {
				diff = json.RawMessage(bodyCopy)
			}
			log := &AuditLog{
				ID:         ulid.Make().String(),
				UserID:     userID,
				EntityName: entity,
				EntityID:   entID,
				Action:     action,
				Diff:       diff,
				IPAddress:  ip,
				UserAgent:  ua,
				GeoInfo:    geoBytes,
				CreatedAt:  time.Now(),
			}
			_ = repo.Create(r.Context(), log)
		})
	}
}
