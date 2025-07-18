package finance

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rgomids/bckoffice/internal/auth"
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

// @Summary      Lista contas a receber
// @Tags         finance
// @Security     BearerAuth
// @Success      200  {array}  AccountReceivable
// @Router       /receivables [get]
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

// @Summary      Marca receivable como pago
// @Tags         finance
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /receivables/{id}/pay [put]
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
	w.Header().Set("X-Entity", fmt.Sprintf("receivables:%s", id))
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Lista comissoes
// @Tags         finance
// @Security     BearerAuth
// @Success      200  {array}  Commission
// @Router       /commissions [get]
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

// @Summary      Aprova comissao
// @Tags         finance
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /commissions/{id}/approve [put]
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
	w.Header().Set("X-Entity", fmt.Sprintf("commissions:%s", id))
	w.WriteHeader(http.StatusNoContent)
}
