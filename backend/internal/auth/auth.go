package auth

import (
	"context"

	"github.com/smithl4b/rcm.backoffice/internal/users"
)

// Repository define operações para consulta de usuários.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (users.User, error)
}
