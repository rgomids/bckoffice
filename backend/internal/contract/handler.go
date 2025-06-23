package contract

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

// RegisterRoutes adiciona as rotas do modulo Contract.
func RegisterRoutes(r chi.Router, repo Repository) {
	h := handler{repo: repo, validate: validator.New()}
	r.Get("/contracts", h.list)
	r.Post("/contracts", h.create)
}

type handler struct {
	repo     Repository
	validate *validator.Validate
}

type createContractInput struct {
	CustomerID string  `json:"customer_id" validate:"required"`
	ServiceID  string  `json:"service_id" validate:"required"`
	PromoterID string  `json:"promoter_id"`
	ValueTotal float64 `json:"value_total" validate:"required,gte=0"`
	StartDate  string  `json:"start_date" validate:"required"`
	EndDate    string  `json:"end_date"`
	Status     string  `json:"status"`
}

func (h handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contracts, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// filtro opcional por status
	status := r.URL.Query().Get("status")
	if status != "" {
		filtered := make([]Contract, 0, len(contracts))
		for _, c := range contracts {
			if c.Status == status {
				filtered = append(filtered, c)
			}
		}
		contracts = filtered
	}

	_ = json.NewEncoder(w).Encode(contracts)
}

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in createContractInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var endDatePtr *time.Time
	if in.EndDate != "" {
		t, err := time.Parse("2006-01-02", in.EndDate)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if startDate.After(t) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "start_date must be before end_date"})
			return
		}
		endDatePtr = &t
	}

	c := Contract{
		ID:         ulid.Make().String(),
		CustomerID: in.CustomerID,
		ServiceID:  in.ServiceID,
	}
	if in.PromoterID != "" {
		c.PromoterID = &in.PromoterID
	}
	c.ValueTotal = in.ValueTotal
	c.StartDate = startDate
	c.EndDate = endDatePtr
	if in.Status != "" {
		c.Status = in.Status
	} else {
		c.Status = "active"
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	if err := h.repo.Create(r.Context(), &c); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/contracts/"+c.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}
