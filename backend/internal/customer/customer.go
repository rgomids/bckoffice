package customer

import (
	"context"
	"errors"
	"time"
)

// Customer representa um cliente da aplicacao.
type Customer struct {
	ID         string     `db:"id" json:"id"`
	LegalName  string     `db:"legal_name" json:"legalName"`
	TradeName  string     `db:"trade_name" json:"tradeName"`
	DocumentID string     `db:"document_id" json:"documentID"`
	Email      string     `db:"email" json:"email"`
	Phone      string     `db:"phone" json:"phone"`
	PromoterID *string    `db:"promoter_id" json:"promoterID,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// Repository define operacoes de acesso ao armazenamento de clientes.
// ErrDuplicateDocumentID eh retornado quando ja existe um cliente com o mesmo
// document_id no banco de dados.
var ErrDuplicateDocumentID = errors.New("duplicate document_id")

type Repository interface {
	FindAll(ctx context.Context) ([]Customer, error)
	FindByID(ctx context.Context, id string) (Customer, error)
	Create(ctx context.Context, c *Customer, addresses []Address) error
	Update(ctx context.Context, c *Customer, addresses []Address) error
	SoftDelete(ctx context.Context, id string) error
}
