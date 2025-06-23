package service

import (
	"encoding/json"
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

func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	services, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(services)
}

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
