package user

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/user/presenter"
)

type UseCase interface {
	GetListUser(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[presenter.User], error)
	DisableUser(ctx context.Context, userId string) error
}
