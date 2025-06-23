package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes adiciona as rotas do módulo Customer.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo}
	r.Get("/customers", h.list)
}

type handler struct {
	repo Repository
}

func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customers, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(customers)
}
