package service

import "time"

// Service representa um servico comercializado pela aplicacao.
type Service struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description,omitempty"`
	BasePrice   float64    `db:"base_price" json:"basePrice"`
	IsActive    bool       `db:"is_active" json:"isActive"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
