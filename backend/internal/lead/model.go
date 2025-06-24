package lead

import "time"

// Lead representa uma oportunidade comercial no pipeline de vendas.
type Lead struct {
	ID         string     `db:"id" json:"id"`
	CustomerID string     `db:"customer_id" json:"customerID"`
	PromoterID *string    `db:"promoter_id" json:"promoterID,omitempty"`
	ServiceID  string     `db:"service_id" json:"serviceID"`
	Status     string     `db:"status" json:"status"`
	Notes      string     `db:"notes" json:"notes"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
