package auditquery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rgomids/bckoffice/internal/auth"
)

// RegisterRoutes adiciona a rota de consulta de audit logs.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo}
	r.Route("/audit-logs", func(rt chi.Router) {
		rt.Use(auth.RequireRole("admin"))
		rt.Get("/", h.list)
	})
}

type handler struct {
	repo Repository
}

// @Summary Lista audit logs
// @Tags audit
// @Security BearerAuth
// @Param entity query string false "Nome da entidade"
// @Param action query string false "insert|update|delete"
// @Success 200 {array} audit.AuditLog
// @Router  /audit-logs [get]
func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	startStr := q.Get("start")
	endStr := q.Get("end")
	var start, end time.Time
	if startStr != "" {
		start, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr != "" {
		end, _ = time.Parse(time.RFC3339, endStr)
	}
	filter := AuditFilter{
		EntityName: q.Get("entity"),
		UserID:     q.Get("user"),
		Action:     q.Get("action"),
		StartDate:  start,
		EndDate:    end,
		Limit:      limit,
	}
	logs, err := h.repo.List(r.Context(), filter)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(logs)
}
