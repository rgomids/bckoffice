package auditquery

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rgomids/bckoffice/internal/audit"
)

// PostgresRepository implementa Repository usando PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository cria uma instancia de PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func toNull(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// List retorna os logs conforme filtros informados.
func (r *PostgresRepository) List(ctx context.Context, f AuditFilter) ([]audit.AuditLog, error) {
	logs := []audit.AuditLog{}
	const q = `SELECT * FROM audit_logs
        WHERE (entity_name=$1 OR $1 IS NULL)
          AND (user_id=$2 OR $2 IS NULL)
          AND (action=$3 OR $3 IS NULL)
          AND created_at BETWEEN $4 AND $5
        ORDER BY created_at DESC
        LIMIT $6`
	start := f.StartDate
	if start.IsZero() {
		start = time.Unix(0, 0)
	}
	end := f.EndDate
	if end.IsZero() {
		end = time.Now()
	}
	if f.Limit == 0 {
		f.Limit = 100
	}
	if err := r.db.SelectContext(ctx, &logs, q,
		toNull(f.EntityName), toNull(f.UserID), toNull(f.Action),
		start, end, f.Limit); err != nil {
		return nil, err
	}
	return logs, nil
}

var _ Repository = (*PostgresRepository)(nil)
