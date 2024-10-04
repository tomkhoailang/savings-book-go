package role

import (
	"context"

	"SavingBooks/internal/domain"
)

type RoleRepository interface {
	domain.GenericRepository[domain.Role]
	TestMethod(ctx context.Context, id string) (*domain.Role, error)
}
