package lead

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

// PostgresRepository implementa Repository usando PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository cria uma instancia de PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// List retorna todos os leads filtrando opcionalmente por status.
func (r *PostgresRepository) List(ctx context.Context, status string) ([]Lead, error) {
	leads := []Lead{}
	q := `SELECT * FROM leads WHERE deleted_at IS NULL`
	if status != "" {
		q += ` AND status = $1`
		if err := r.db.SelectContext(ctx, &leads, q, status); err != nil {
			return nil, err
		}
	} else {
		if err := r.db.SelectContext(ctx, &leads, q); err != nil {
			return nil, err
		}
	}
	return leads, nil
}

// Create insere um novo lead gerando ULID.
func (r *PostgresRepository) Create(ctx context.Context, l *Lead) error {
	l.ID = ulid.Make().String()
	const q = `INSERT INTO leads (id, customer_id, promoter_id, service_id, status, notes)
        VALUES (:id, :customer_id, :promoter_id, :service_id, :status, :notes)`
	_, err := r.db.NamedExecContext(ctx, q, l)
	return err
}

var allowedTransitions = map[string]string{
	"lead":      "qualified",
	"qualified": "proposal",
	"proposal":  "contract",
}

// UpdateStatus atualiza o status de um lead verificando transicao valida.
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	var current string
	if err := r.db.GetContext(ctx, &current, `SELECT status FROM leads WHERE id=$1 AND deleted_at IS NULL`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	if allowedTransitions[current] != newStatus {
		return errors.New("invalid status transition")
	}
	const q = `UPDATE leads SET status=$2, updated_at=now() WHERE id=$1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id, newStatus)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Update altera dados de um lead.
func (r *PostgresRepository) Update(ctx context.Context, l *Lead) error {
	const q = `UPDATE leads SET promoter_id=:promoter_id, service_id=:service_id, notes=:notes, updated_at=now()
        WHERE id=:id AND deleted_at IS NULL`
	res, err := r.db.NamedExecContext(ctx, q, l)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// SoftDelete marca um lead como removido.
func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE leads SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
