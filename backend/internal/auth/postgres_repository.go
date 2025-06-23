package auth

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/smithl4b/rcm.backoffice/internal/users"
)

// PostgresRepository implementa Repository usando PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository cria uma instancia de PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// FindByEmail retorna um usuario pelo e-mail.
func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (users.User, error) {
	var u users.User
	const q = `SELECT u.id, u.email, u.password_hash, u.full_name, u.created_at, u.updated_at, u.deleted_at,
        COALESCE(r.name,'') AS role
        FROM users u
        LEFT JOIN user_roles ur ON ur.user_id = u.id
        LEFT JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL
        WHERE u.email=$1 AND u.deleted_at IS NULL LIMIT 1`
	if err := r.db.GetContext(ctx, &u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return users.User{}, sql.ErrNoRows
		}
		return users.User{}, err
	}
	return u, nil
}
