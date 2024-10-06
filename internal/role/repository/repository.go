package repository

import (
	"context"
	"fmt"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleRepository struct {
	contracts.BaseRepository[domain.Role]
}

func (r *roleRepository) SeedRole(ctx context.Context) error {

	count, err := r.CountAll(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		print("No need to seed role")
		return nil
	}

	var input = make([]domain.Role, 2)
	input[0] = domain.Role{Name: "Admin"}
	input[1] = domain.Role{Name: "User"}
	for _, role := range input {
		err := r.Create(ctx, &role)
		if err != nil {
			return err
		}
	}
	fmt.Println("Seeding roles completed")
	return nil
}

func NewRoleRepository(db *mongo.Database, collectionName string) role.RoleRepository {
	baseRepo := contracts.NewBaseRepository[domain.Role](db, collectionName).(*contracts.BaseRepository[domain.Role])
	return &roleRepository{BaseRepository: *baseRepo }
}


