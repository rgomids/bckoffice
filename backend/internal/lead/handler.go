package lead

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"

	"github.com/rgomids/bckoffice/internal/auth"
)

// RegisterRoutes adiciona as rotas do modulo Lead.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo, validate: validator.New()}

	r.Group(func(gr chi.Router) {
		gr.Use(auth.RequireRole("admin", "promoter"))
		gr.Get("/leads", h.list)
	})

	r.Post("/leads", h.create)
	r.With(auth.RequireRole("finance", "admin")).Put("/leads/{id}/status", h.updateStatus)
	r.Put("/leads/{id}", h.update)
	r.Delete("/leads/{id}", h.remove)
}

type handler struct {
	repo     Repository
	validate *validator.Validate
}

type createLeadInput struct {
	CustomerID string  `json:"customer_id" validate:"required"`
	ServiceID  string  `json:"service_id" validate:"required"`
	PromoterID *string `json:"promoter_id"`
	Notes      string  `json:"notes"`
}

type updateLeadInput struct {
	ServiceID  string  `json:"service_id" validate:"required"`
	PromoterID *string `json:"promoter_id"`
	Notes      string  `json:"notes"`
}

type statusInput struct {
	Status string `json:"status" validate:"required,oneof=qualified proposal contract"`
}

// @Summary      Lista leads
// @Tags         leads
// @Security     BearerAuth
// @Success      200  {array}  Lead
// @Router       /leads [get]
func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := r.URL.Query().Get("status")
	leads, err := h.repo.List(r.Context(), status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(leads)
}

// @Summary      Cria lead
// @Tags         leads
// @Security     BearerAuth
// @Success      201  {object}  Lead
// @Router       /leads [post]
func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in createLeadInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	l := Lead{
		ID:         ulid.Make().String(),
		CustomerID: in.CustomerID,
		ServiceID:  in.ServiceID,
		PromoterID: in.PromoterID,
		Status:     "lead",
		Notes:      in.Notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.repo.Create(r.Context(), &l); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/leads/"+l.ID)
	w.Header().Set("X-Entity", fmt.Sprintf("leads:%s", l.ID))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(l)
}

// @Summary      Atualiza status do lead
// @Tags         leads
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /leads/{id}/status [put]
func (h handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var in statusInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, in.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err.Error() == "invalid status transition" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Entity", fmt.Sprintf("leads:%s", id))
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Atualiza lead
// @Tags         leads
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /leads/{id} [put]
func (h handler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	var in updateLeadInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	l := Lead{
		ID:         id,
		ServiceID:  in.ServiceID,
		PromoterID: in.PromoterID,
		Notes:      in.Notes,
		UpdatedAt:  time.Now(),
	}

	if err := h.repo.Update(r.Context(), &l); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Entity", fmt.Sprintf("leads:%s", id))
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Remove lead
// @Tags         leads
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /leads/{id} [delete]
func (h handler) remove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.SoftDelete(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-Entity", fmt.Sprintf("leads:%s", id))
	w.WriteHeader(http.StatusNoContent)
}
