package promoter

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
