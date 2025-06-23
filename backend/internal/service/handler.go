package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

// RegisterRoutes adiciona as rotas do modulo Service.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo, validate: validator.New()}
	r.Get("/services", h.list)
	r.Post("/services", h.create)
	r.Put("/services/{id}", h.update)
	r.Delete("/services/{id}", h.remove)
}

type handler struct {
	repo     Repository
	validate *validator.Validate
}

type createServiceInput struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" validate:"gte=0"`
	IsActive    bool    `json:"is_active"`
}

// UpdateServiceInput define o payload para atualizacao de servicos.
type UpdateServiceInput struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" validate:"gte=0"`
	IsActive    bool    `json:"is_active"`
}

// @Summary      Lista servicos
// @Tags         services
// @Security     BearerAuth
// @Success      200  {array}  Service
// @Router       /services [get]
func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	services, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(services)
}

// @Summary      Cria servico
// @Tags         services
// @Security     BearerAuth
// @Success      201  {object}  Service
// @Router       /services [post]
func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in createServiceInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	s := Service{
		ID:          ulid.Make().String(),
		Name:        in.Name,
		Description: in.Description,
		BasePrice:   in.BasePrice,
		IsActive:    in.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(r.Context(), &s); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/services/"+s.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

// @Summary      Atualiza servico
// @Tags         services
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /services/{id} [put]
func (h handler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	var in UpdateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	s := Service{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		BasePrice:   in.BasePrice,
		IsActive:    in.IsActive,
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Update(r.Context(), &s); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Remove servico
// @Tags         services
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /services/{id} [delete]
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
	w.WriteHeader(http.StatusNoContent)
}
