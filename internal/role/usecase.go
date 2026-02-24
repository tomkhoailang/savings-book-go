package role

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role/presenter"
)

type UseCase interface {
	CreateRole(ctx context.Context, input *presenter.RoleInput, creatorId string) (*domain.Role, error)
	UpdateRole(ctx context.Context, input *presenter.RoleInput, lastModifierId, roleId string) (*domain.Role, error)
	DeleteManyRoles(ctx context.Context, deleterId string, input []string) error
	GetListRoles(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[domain.Role], error)
	SeedRoles(ctx context.Context) error
}