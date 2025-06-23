package finance

import (
	"context"
	"errors"
)

// Repository define operacoes para contas a receber e comissoes.
type Repository interface {
	ListReceivables(ctx context.Context, status string) ([]AccountReceivable, error)
	MarkAsPaid(ctx context.Context, id string) error

	ListCommissions(ctx context.Context, onlyPending bool) ([]Commission, error)
	ApproveCommission(ctx context.Context, id string, approverID string) error
}

// Errors especificos

var (
	ErrAlreadyPaid     = errors.New("already paid")
	ErrAlreadyApproved = errors.New("already approved")
)
