package promoter

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

// RegisterRoutes adiciona as rotas do módulo Promoter.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo, validate: validator.New()}
	r.Get("/promoters", h.list)
	r.Post("/promoters", h.create)
}

type handler struct {
	repo     Repository
	validate *validator.Validate
}

type createPromoterInput struct {
	FullName    string          `json:"full_name" validate:"required"`
	Email       string          `json:"email" validate:"omitempty,email"`
	Phone       string          `json:"phone"`
	DocumentID  string          `json:"document_id"`
	BankAccount json.RawMessage `json:"bank_account"`
}

func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	promoters, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(promoters)
}

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in createPromoterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	p := Promoter{
		ID:          ulid.Make().String(),
		FullName:    in.FullName,
		Email:       in.Email,
		Phone:       in.Phone,
		DocumentID:  in.DocumentID,
		BankAccount: in.BankAccount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(r.Context(), &p); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/promoters/"+p.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}
