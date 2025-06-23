package contract

import "time"

// Contract representa um contrato entre cliente e serviço.
type Contract struct {
	ID         string     `db:"id" json:"id"`
	CustomerID string     `db:"customer_id" json:"customer_id"`
	ServiceID  string     `db:"service_id" json:"service_id"`
	PromoterID *string    `db:"promoter_id" json:"promoter_id,omitempty"`
	ValueTotal float64    `db:"value_total" json:"value_total"`
	StartDate  time.Time  `db:"start_date" json:"start_date"`
	EndDate    *time.Time `db:"end_date" json:"end_date,omitempty"`
	Status     string     `db:"status" json:"status"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
