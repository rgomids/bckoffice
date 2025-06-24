package auditquery

import (
	"context"
	"time"

	"github.com/rgomids/bckoffice/internal/audit"
)

// AuditFilter define filtros para listagem de logs.
type AuditFilter struct {
	EntityName string
	UserID     string
	Action     string
	StartDate  time.Time
	EndDate    time.Time
	Limit      int
}

// Repository define operacoes de consulta aos logs de auditoria.
type Repository interface {
	List(ctx context.Context, filter AuditFilter) ([]audit.AuditLog, error)
}
