package auth

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error)
	SignIn(ctx context.Context, creds presenter.LoginInput) (*presenter.LogInRes, error)
	ParseAccessToken(accessToken string) (*presenter.TokenResult, error)
	RenewAccessToken(ctx context.Context, req *presenter.RenewTokenReq) (string, error)
	Logout(ctx context.Context, userId string) error
}
