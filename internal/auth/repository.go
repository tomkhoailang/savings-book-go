package auth

import (
	"context"

	"SavingBooks/internal/domain"
)

type UserRepository interface {
	domain.GenericRepository[domain.User]
	GetExistUser(ctx context.Context, username, email string) (*domain.User, error)
}
