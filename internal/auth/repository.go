package auth

import (
	"context"

	"SavingBooks/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
}
