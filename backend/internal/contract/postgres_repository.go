package contract

import (
	"context"
	"database/sql"
	"errors"

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

// FindAll retorna todos os contratos nao excluidos.
func (r *PostgresRepository) FindAll(ctx context.Context) ([]Contract, error) {
	contracts := []Contract{}
	const q = `SELECT * FROM contracts WHERE deleted_at IS NULL ORDER BY start_date DESC`
	if err := r.db.SelectContext(ctx, &contracts, q); err != nil {
		return nil, err
	}
	return contracts, nil
}

// Create insere um novo contrato.
func (r *PostgresRepository) Create(ctx context.Context, c *Contract) error {
	// valida existencia de customer e service
	var exists int
	if err := r.db.GetContext(ctx, &exists, `SELECT 1 FROM customers WHERE id=$1 AND deleted_at IS NULL`, c.CustomerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	if err := r.db.GetContext(ctx, &exists, `SELECT 1 FROM services WHERE id=$1 AND deleted_at IS NULL`, c.ServiceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	const q = `INSERT INTO contracts (id, customer_id, service_id, promoter_id, value_total, start_date, end_date, status) VALUES (:id, :customer_id, :service_id, :promoter_id, :value_total, :start_date, :end_date, :status)`
	_, err := r.db.NamedExecContext(ctx, q, c)
	return err
}
