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

// Update altera dados de um contrato existente.
func (r *PostgresRepository) Update(ctx context.Context, c *Contract) error {
       const q = `UPDATE contracts SET value_total=:value_total, start_date=:start_date, end_date=:end_date, status=:status, updated_at=now() WHERE id=:id AND deleted_at IS NULL`
       res, err := r.db.NamedExecContext(ctx, q, c)
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

// SoftDelete marca um contrato como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
       const q = `UPDATE contracts SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
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
