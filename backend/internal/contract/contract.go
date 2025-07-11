package contract

import "context"

// Repository define operações para persistência de contratos.
type Repository interface {
	FindAll(ctx context.Context) ([]Contract, error)
	Create(ctx context.Context, c *Contract) error
	Update(ctx context.Context, c *Contract) error
	SoftDelete(ctx context.Context, id string) error
}
