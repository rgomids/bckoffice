package audit

import "context"

// Repository define operações para persistir logs de auditoria.
type Repository interface {
	Create(ctx context.Context, log *AuditLog) error
}
