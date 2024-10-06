package usecase

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"SavingBooks/internal/role/presenter"
	"github.com/jinzhu/copier"
)

type roleUseCase struct {
	roleRepo role.RoleRepository
}


func (r *roleUseCase) GetListRoles(ctx context.Context, query *contracts.Query) (contracts.QueryResult[domain.Role], error) {
	rolesInterface, err := r.roleRepo.GetList(ctx, query)
	if err != nil {
		return contracts.QueryResult[domain.Role]{}, err
	}
	roles := rolesInterface.(*contracts.QueryResult[domain.Role])

	return *roles, nil
}

func (r *roleUseCase) CreateRole(ctx context.Context, input *presenter.RoleInput, creatorId string) (*domain.Role, error) {
	role := &domain.Role{
		Name:        input.Name,
		Description: input.Description,
	}
	role.SetCreate(creatorId)
	role.Keyword = input.Name + " " + input.Description
	err := r.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleUseCase) UpdateRole(ctx context.Context, input *presenter.RoleInput, lastModifierId, roleId string) (*domain.Role, error) {
	entity, err := r.roleRepo.Get(ctx, roleId)
	if err != nil{
		return nil, role.ErrRoleNotFound
	}
	err = copier.Copy(entity, input)
	if err != nil {
		return nil, err
	}
	entity.Keyword = input.Name + " " + input.Description
	entity.SetUpdate(lastModifierId)

	err = r.roleRepo.Update(ctx, entity, roleId)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *roleUseCase) DeleteManyRoles(ctx context.Context, deleterId string, input []string) error {
	err := r.roleRepo.DeleteMany(ctx, deleterId, input)
	if err != nil {
		return err
	}
	return nil
}

func (r *roleUseCase) SeedRoles(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewRoleUseCase(roleRepo role.RoleRepository) role.UseCase {
	return &roleUseCase{roleRepo: roleRepo}
}


