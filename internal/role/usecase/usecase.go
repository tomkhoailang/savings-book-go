package usecase

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"SavingBooks/internal/role/presenter"
)

type roleUseCase struct {
	roleRepo role.RoleRepository
}

func (r *roleUseCase) GetListRoles(ctx context.Context, query contracts.Query) ([]domain.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleUseCase) CreateRole(ctx context.Context, input *presenter.RoleInput, creatorId string) (*domain.Role, error) {
	role := &domain.Role{
		Name:        input.Name,
		Description: input.Name,
	}
	role.SetCreate(creatorId)
	role.Keyword = input.Name + " " + input.Description
	err := r.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleUseCase) UpdateRole(ctx context.Context, role *presenter.RoleInput, lastModifierId string) (*domain.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleUseCase) DeleteManyRoles(ctx context.Context, role *presenter.RoleInput, deleterId string) error {
	//TODO implement me
	panic("implement me")
}

func NewRoleUseCase(roleRepo role.RoleRepository) role.UseCase {
	return &roleUseCase{roleRepo: roleRepo}
}


