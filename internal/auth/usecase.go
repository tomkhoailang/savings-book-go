package auth

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error)
	SignIn(ctx context.Context, creds presenter.LoginInput) (string, error)
	ParseAccessToken(ctx context.Context, accessToken string) (*presenter.TokenResult, error)
}
