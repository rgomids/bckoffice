package service

import (
	"context"
	"database/sql"

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

// Update atualiza um servico existente.
func (r *PostgresRepository) Update(ctx context.Context, s *Service) error {
	const q = `UPDATE services SET name=:name, description=:description, base_price=:base_price, is_active=:is_active, updated_at=now() WHERE id=:id AND deleted_at IS NULL`
	res, err := r.db.NamedExecContext(ctx, q, s)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// SoftDelete marca um servico como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	const q = `UPDATE services SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
