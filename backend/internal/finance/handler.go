package finance

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smithl4b/rcm.backoffice/internal/auth"
)

// RegisterRoutes adiciona as rotas do modulo Finance.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo}
	r.Route("/receivables", func(r chi.Router) {
		r.Use(auth.RequireRole("finance"))
		r.Get("/", h.listReceivables)
		r.Put("/{id}/pay", h.markAsPaid)
	})
	r.Route("/commissions", func(r chi.Router) {
		r.Use(auth.RequireRole("finance"))
		r.Get("/", h.listCommissions)
		r.Put("/{id}/approve", h.approveCommission)
	})
}

type handler struct {
	repo Repository
}

func (h handler) listReceivables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := r.URL.Query().Get("status")
	list, err := h.repo.ListReceivables(r.Context(), status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h handler) markAsPaid(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.repo.MarkAsPaid(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrAlreadyPaid) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) listCommissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pending := r.URL.Query().Get("pending") == "true"
	list, err := h.repo.ListCommissions(r.Context(), pending)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h handler) approveCommission(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := auth.UserIDFromContext(r.Context())
	err := h.repo.ApproveCommission(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrAlreadyApproved) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
