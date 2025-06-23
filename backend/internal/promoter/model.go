package promoter

import (
	"encoding/json"
	"time"
)

// Promoter representa um divulgador de serviços.
type Promoter struct {
	ID          string          `db:"id" json:"id"`
	FullName    string          `db:"full_name" json:"fullName"`
	Email       string          `db:"email" json:"email,omitempty"`
	Phone       string          `db:"phone" json:"phone,omitempty"`
	DocumentID  string          `db:"document_id" json:"documentID,omitempty"`
	BankAccount json.RawMessage `db:"bank_account" json:"bankAccount,omitempty"`
	CreatedAt   time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time      `db:"deleted_at" json:"deletedAt,omitempty"`
}
