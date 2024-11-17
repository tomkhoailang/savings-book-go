package auth

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error)
	SignIn(ctx context.Context, creds presenter.LoginInput) (*presenter.LogInRes, error)
	GenerateResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error
	ConfirmResetPassword(ctx context.Context, token, password string) error
	ParseAccessToken(ctx context.Context, accessToken string) (*presenter.TokenResult, error)
	RenewAccessToken(ctx context.Context, req *presenter.RenewTokenReq) (string, error)
	Logout(ctx context.Context, userId string) error
}
