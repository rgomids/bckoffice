package customer

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
)

// RegisterRoutes adiciona as rotas do módulo Customer.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo, validate: validator.New()}
	r.Get("/customers", h.list)
	r.Post("/customers", h.create)
	r.Put("/customers/{id}", h.update)
	r.Delete("/customers/{id}", h.remove)
}

type handler struct {
	repo     Repository
	validate *validator.Validate
}

// @Summary      Lista clientes
// @Tags         customers
// @Security     BearerAuth
// @Success      200  {array}  Customer
// @Router       /customers [get]
func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customers, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(customers)
}

// @Summary      Cria cliente
// @Tags         customers
// @Security     BearerAuth
// @Success      201  {object}  Customer
// @Router       /customers [post]
func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in CreateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	c := Customer{
		ID:         ulid.Make().String(),
		LegalName:  in.LegalName,
		TradeName:  in.TradeName,
		DocumentID: in.DocumentID,
		Email:      in.Email,
		Phone:      in.Phone,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	addresses := make([]Address, len(in.Addresses))
	for i, a := range in.Addresses {
		addresses[i] = Address{
			ID:          ulid.Make().String(),
			CustomerID:  c.ID,
			AddressType: a.AddressType,
			Street:      a.Street,
			Number:      a.Number,
			Complement:  a.Complement,
			District:    a.District,
			City:        a.City,
			State:       a.State,
			PostalCode:  a.PostalCode,
			Country:     a.Country,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	if err := h.repo.Create(r.Context(), &c, addresses); err != nil {
		if errors.Is(err, ErrDuplicateDocumentID) {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/customers/"+c.ID)
	w.Header().Set("X-Entity", fmt.Sprintf("customers:%s", c.ID))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// @Summary      Atualiza cliente
// @Tags         customers
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /customers/{id} [put]
func (h handler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	var in UpdateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	c := Customer{
		ID:         id,
		LegalName:  in.LegalName,
		TradeName:  in.TradeName,
		DocumentID: in.DocumentID,
		Email:      in.Email,
		Phone:      in.Phone,
		UpdatedAt:  time.Now(),
	}

	var addresses []Address
	if in.Addresses != nil {
		addresses = make([]Address, len(in.Addresses))
		for i, a := range in.Addresses {
			addresses[i] = Address{
				ID:          ulid.Make().String(),
				CustomerID:  id,
				AddressType: a.AddressType,
				Street:      a.Street,
				Number:      a.Number,
				Complement:  a.Complement,
				District:    a.District,
				City:        a.City,
				State:       a.State,
				PostalCode:  a.PostalCode,
				Country:     a.Country,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
		}
	}

	if err := h.repo.Update(r.Context(), &c, addresses); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrDuplicateDocumentID) {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Entity", fmt.Sprintf("customers:%s", id))
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Remove cliente
// @Tags         customers
// @Security     BearerAuth
// @Success      204  {null}  nil
// @Router       /customers/{id} [delete]
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
	w.Header().Set("X-Entity", fmt.Sprintf("customers:%s", id))
	w.WriteHeader(http.StatusNoContent)
}
