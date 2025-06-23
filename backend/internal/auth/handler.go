package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRoutes adiciona a rota de login.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo}
	r.Post("/login", h.login)
}

type handler struct {
	repo Repository
}

// @Summary      Autentica usuario
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      CredentialsInput  true  "Credenciais de login"
// @Success      200  {object}  AuthResponse
// @Router       /login [post]
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in CredentialsInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindByEmail(r.Context(), in.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(AuthResponse{Token: token})
}
