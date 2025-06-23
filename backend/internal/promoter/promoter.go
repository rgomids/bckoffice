package promoter

import "context"

// Repository define operações para armazenamento de promotores.
type Repository interface {
	FindAll(ctx context.Context) ([]Promoter, error)
	Create(ctx context.Context, p *Promoter) error
}
