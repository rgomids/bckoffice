package finance

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

// ListReceivables retorna as contas a receber opcionamente filtrando por status.
func (r *PostgresRepository) ListReceivables(ctx context.Context, status string) ([]AccountReceivable, error) {
	receivables := []AccountReceivable{}
	q := `SELECT * FROM accounts_receivable WHERE deleted_at IS NULL`
	if status != "" {
		q += ` AND status=$1`
		if err := r.db.SelectContext(ctx, &receivables, q, status); err != nil {
			return nil, err
		}
	} else {
		if err := r.db.SelectContext(ctx, &receivables, q); err != nil {
			return nil, err
		}
	}
	return receivables, nil
}

// MarkAsPaid marca uma conta como paga.
func (r *PostgresRepository) MarkAsPaid(ctx context.Context, id string) error {
	const q = `UPDATE accounts_receivable SET status='paid', paid_at=now() WHERE id=$1 AND status='open'`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 1 {
		return nil
	}
	var status string
	err = r.db.GetContext(ctx, &status, `SELECT status FROM accounts_receivable WHERE id=$1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	return ErrAlreadyPaid
}

// ListCommissions retorna as comissoes, opcionalmente apenas pendentes.
func (r *PostgresRepository) ListCommissions(ctx context.Context, onlyPending bool) ([]Commission, error) {
	commissions := []Commission{}
	q := `SELECT * FROM commissions WHERE deleted_at IS NULL`
	if onlyPending {
		q += ` AND approved=false`
	}
	if err := r.db.SelectContext(ctx, &commissions, q); err != nil {
		return nil, err
	}
	return commissions, nil
}

// ApproveCommission marca uma comissao como aprovada.
func (r *PostgresRepository) ApproveCommission(ctx context.Context, id string, approverID string) error {
	const q = `UPDATE commissions SET approved=true, approved_by=$2, approved_at=now() WHERE id=$1 AND approved=false`
	res, err := r.db.ExecContext(ctx, q, id, approverID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 1 {
		return nil
	}
	var approved bool
	err = r.db.GetContext(ctx, &approved, `SELECT approved FROM commissions WHERE id=$1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	return ErrAlreadyApproved
}

var _ Repository = (*PostgresRepository)(nil)
