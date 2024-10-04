package role

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role/presenter"
)

type UseCase interface {
	CreateRole(ctx context.Context, input *presenter.RoleInput, creatorId string) (*domain.Role, error)
	UpdateRole(ctx context.Context, role *presenter.RoleInput, lastModifierId string) (*domain.Role, error)
	DeleteManyRoles(ctx context.Context, role *presenter.RoleInput, deleterId string) error
	GetListRoles(ctx context.Context, query contracts.Query) ([]domain.Role, error)
}