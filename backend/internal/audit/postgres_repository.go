package audit

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

// Create insere um novo registro de auditoria.
func (r *PostgresRepository) Create(ctx context.Context, log *AuditLog) error {
	const q = `INSERT INTO audit_logs (
        id, user_id, entity_name, entity_id, action, diff, ip_address,
        user_agent, geo_info)
        VALUES (:id, :user_id, :entity_name, :entity_id, :action, :diff,
                :ip_address, :user_agent, :geo_info)`
	_, err := r.db.NamedExecContext(ctx, q, log)
	return err
}
