package user

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	GetListUser(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[domain.User], error)
	DisableUser(ctx context.Context, userId string) error
}
