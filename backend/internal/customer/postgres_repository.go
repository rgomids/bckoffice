package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PostgresRepository implementa Repository usando PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository cria uma instancia de PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// FindAll retorna todos os clientes nao excluidos.
func (r *PostgresRepository) FindAll(ctx context.Context) ([]Customer, error) {
	customers := []Customer{}
	const q = `SELECT * FROM customers WHERE deleted_at IS NULL`
	if err := r.db.SelectContext(ctx, &customers, q); err != nil {
		return nil, err
	}
	return customers, nil
}

// FindByID retorna um cliente pelo ID.
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (Customer, error) {
	var c Customer
	const q = `SELECT * FROM customers WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &c, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Customer{}, nil
		}
		return Customer{}, err
	}
	return c, nil
}

// Create insere um novo cliente.
func (r *PostgresRepository) Create(ctx context.Context, c *Customer, addresses []Address) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	const qc = `INSERT INTO customers (id, legal_name, trade_name, document_id, email, phone, promoter_id)
                VALUES (:id, :legal_name, :trade_name, :document_id, :email, :phone, :promoter_id)`
	if _, err = tx.NamedExecContext(ctx, qc, c); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			_ = tx.Rollback()
			return ErrDuplicateDocumentID
		}
		_ = tx.Rollback()
		return err
	}

	const qa = `INSERT INTO addresses (id, customer_id, address_type, street, number, complement, district, city, state, postal_code, country)
                VALUES (:id, :customer_id, :address_type, :street, :number, :complement, :district, :city, :state, :postal_code, :country)`
	for _, a := range addresses {
		if _, err = tx.NamedExecContext(ctx, qa, a); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert address: %w", err)
		}
	}

	return tx.Commit()
}

// Update atualiza dados do cliente.
func (r *PostgresRepository) Update(ctx context.Context, c *Customer, addrs []Address) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	const qc = `UPDATE customers SET legal_name=:legal_name, trade_name=:trade_name, document_id=:document_id,
                email=:email, phone=:phone, promoter_id=:promoter_id, updated_at=now()
                WHERE id=:id AND deleted_at IS NULL`
	res, err := tx.NamedExecContext(ctx, qc, c)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			_ = tx.Rollback()
			return ErrDuplicateDocumentID
		}
		_ = tx.Rollback()
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if affected == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}

	if addrs != nil {
		if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE customer_id=$1`, c.ID); err != nil {
			_ = tx.Rollback()
			return err
		}

		const qa = `INSERT INTO addresses (id, customer_id, address_type, street, number, complement, district, city, state, postal_code, country)
                        VALUES (:id, :customer_id, :address_type, :street, :number, :complement, :district, :city, :state, :postal_code, :country)`
		for _, a := range addrs {
			if _, err = tx.NamedExecContext(ctx, qa, a); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("insert address: %w", err)
			}
		}
	}

	return tx.Commit()
}

// SoftDelete marca um cliente como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	const q = `UPDATE customers SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
