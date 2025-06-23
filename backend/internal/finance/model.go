package finance

import "time"

// AccountReceivable representa um valor a receber de um contrato.
type AccountReceivable struct {
	ID         string     `db:"id" json:"id"`
	ContractID string     `db:"contract_id" json:"contractID"`
	DueDate    time.Time  `db:"due_date" json:"dueDate"`
	Amount     float64    `db:"amount" json:"amount"`
	Status     string     `db:"status" json:"status"`
	PaidAt     *time.Time `db:"paid_at" json:"paidAt,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}

// Commission representa a comissao de um promotor por contrato.
type Commission struct {
	ID         string     `db:"id" json:"id"`
	ContractID string     `db:"contract_id" json:"contractID"`
	PromoterID string     `db:"promoter_id" json:"promoterID"`
	Amount     float64    `db:"amount" json:"amount"`
	Approved   bool       `db:"approved" json:"approved"`
	ApprovedBy string     `db:"approved_by" json:"approvedBy,omitempty"`
	ApprovedAt *time.Time `db:"approved_at" json:"approvedAt,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
