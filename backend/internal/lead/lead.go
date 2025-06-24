package lead

import "context"

// Repository define operacoes para gerenciar leads de vendas.
type Repository interface {
	List(ctx context.Context, statusFilter string) ([]Lead, error)
	Create(ctx context.Context, l *Lead) error
	UpdateStatus(ctx context.Context, id string, newStatus string) error
	Update(ctx context.Context, l *Lead) error
	SoftDelete(ctx context.Context, id string) error
}
