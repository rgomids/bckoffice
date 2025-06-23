package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// PostgresRepository implementa Repository usando PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository cria uma instancia de PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// FindAll retorna todos os servicos ativos (nao excluidos).
func (r *PostgresRepository) FindAll(ctx context.Context) ([]Service, error) {
	services := []Service{}
	const q = `SELECT * FROM services WHERE deleted_at IS NULL ORDER BY name`
	if err := r.db.SelectContext(ctx, &services, q); err != nil {
		return nil, err
	}
	return services, nil
}

// Create insere um novo servico.
func (r *PostgresRepository) Create(ctx context.Context, s *Service) error {
	const q = `INSERT INTO services (id, name, description, base_price, is_active) VALUES (:id, :name, :description, :base_price, :is_active)`
	_, err := r.db.NamedExecContext(ctx, q, s)
	return err
}
