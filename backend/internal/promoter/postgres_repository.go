package promoter

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

// FindAll retorna todos os promotores nao excluidos.
func (r *PostgresRepository) FindAll(ctx context.Context) ([]Promoter, error) {
	promoters := []Promoter{}
	const q = `SELECT * FROM promoters WHERE deleted_at IS NULL ORDER BY full_name`
	if err := r.db.SelectContext(ctx, &promoters, q); err != nil {
		return nil, err
	}
	return promoters, nil
}

// Create insere um novo promotor.
func (r *PostgresRepository) Create(ctx context.Context, p *Promoter) error {
	const q = `INSERT INTO promoters (id, full_name, email, phone, document_id, bank_account) VALUES (:id, :full_name, :email, :phone, :document_id, :bank_account)`
	_, err := r.db.NamedExecContext(ctx, q, p)
	return err
}

// Update atualiza um promotor existente.
func (r *PostgresRepository) Update(ctx context.Context, p *Promoter) error {
	const q = `UPDATE promoters SET full_name=:full_name, email=:email, phone=:phone, document_id=:document_id, bank_account=:bank_account, updated_at=now() WHERE id=:id AND deleted_at IS NULL`
	res, err := r.db.NamedExecContext(ctx, q, p)
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

// SoftDelete marca um promotor como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	const q = `UPDATE promoters SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`
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
