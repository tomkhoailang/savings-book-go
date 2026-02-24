package role

import (
	"context"

	"SavingBooks/internal/domain"
)

type RoleRepository interface {
	domain.GenericRepository[domain.Role]
	SeedRole(ctx context.Context) error
}
