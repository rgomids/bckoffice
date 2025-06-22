package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes adiciona as rotas do módulo Customer.
func RegisterRoutes(r chi.Router) {
	r.Get("/customers", list)
}

func list(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]interface{}{}) // placeholder []
}
