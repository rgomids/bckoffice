package customer

import (
        "encoding/json"
        "errors"
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
}

type handler struct {
        repo     Repository
        validate *validator.Validate
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
        w.WriteHeader(http.StatusCreated)
        _ = json.NewEncoder(w).Encode(c)
}
