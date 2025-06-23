package service

import "context"

// Repository define operacoes de acesso aos servicos.
type Repository interface {
	FindAll(ctx context.Context) ([]Service, error)
	Create(ctx context.Context, s *Service) error
	Update(ctx context.Context, s *Service) error
	SoftDelete(ctx context.Context, id string) error
}
