package customer

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
func (r *PostgresRepository) Create(ctx context.Context, c *Customer) error {
	const q = `INSERT INTO customers (id, legal_name, trade_name, document_id, email, phone, promoter_id)
               VALUES (:id, :legal_name, :trade_name, :document_id, :email, :phone, :promoter_id)`
	_, err := r.db.NamedExecContext(ctx, q, c)
	return err
}

// Update atualiza dados do cliente.
func (r *PostgresRepository) Update(ctx context.Context, c *Customer) error {
	const q = `UPDATE customers SET legal_name=:legal_name, trade_name=:trade_name, document_id=:document_id,
                email=:email, phone=:phone, promoter_id=:promoter_id, updated_at=now()
                WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, q, c)
	return err
}

// SoftDelete marca um cliente como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	const q = `UPDATE customers SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
