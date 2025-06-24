package auth

import (
	"context"

	"github.com/rgomids/bckoffice/internal/users"
)

// Repository define operações para consulta de usuários.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (users.User, error)
}
