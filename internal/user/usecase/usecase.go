package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"SavingBooks/internal/role"
	"SavingBooks/internal/services/redis/redis_key"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	"SavingBooks/internal/services/redis"
	"SavingBooks/internal/user"
	"SavingBooks/internal/user/presenter"
	"github.com/jinzhu/copier"
)

type userUseCase struct {
	userRepo     auth.UserRepository
	roleRepo     role.RoleRepository
	cacheService *redis.Cache
	monthlyRepo  monthly_saving_interest.Repository
	expireTime   time.Duration
}

func (u *userUseCase) GetListUser(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[presenter.User], error) {
	var userInterfaces interface{}
	var err error
	userInterfaces, err = u.userRepo.GetList(ctx, query)

	if err != nil {
		return nil, err
	}

	userList := userInterfaces.(*contracts.QueryResult[domain.User])
	userListOutput := &contracts.QueryResult[presenter.User]{}
	err = copier.Copy(&userListOutput, userList)
	if err != nil {
		return nil, err
	}
	return userListOutput, nil
}

func (u *userUseCase) DisableUser(ctx context.Context, userId string) error {
	user, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}

	collectionInterface := u.roleRepo.GetCollection()
	collection := collectionInterface.(*mongo.Collection)

	query := bson.M{"_id": bson.M{"$in": user.RoleIds}}

	var roles []domain.Role
	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		result := domain.Role{}
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding saving book:", err)
			continue
		}
		roles = append(roles, result)
	}

	if err := cursor.Err(); err != nil {
		return err
	}
	isAdmin := false
	for _, curRole := range roles {
		if curRole.Name == "Admin" {
			isAdmin = true
			break
		}
	}

	if isAdmin {
		return errors.New("Admin couldn't be blocked")
	}

	user.IsActive = !user.IsActive
	if !user.IsActive {
		err = u.cacheService.SetValueWithExpire(ctx, redis_key.BlockUserId+":"+userId, user.Id, 24*time.Hour)
		if err != nil {
			return err
		}
	} else {
		err := u.cacheService.RemoveValue(ctx, redis_key.BlockUserId+":"+userId)
		if err != nil {
			return err
		}
	}
	user.RefreshToken = ""
	_, err = u.userRepo.Update(ctx, user, userId, []string{"IsActive", "RefreshToken"})
	if err != nil {
		return err
	}

	return nil
}

func NewUserUseCase(userRepo auth.UserRepository, monthlyRepo monthly_saving_interest.Repository, roleRepo role.RoleRepository, cacheService *redis.Cache, expireTime time.Duration) user.UseCase {
	return &userUseCase{userRepo: userRepo, cacheService: cacheService, expireTime: expireTime, monthlyRepo: monthlyRepo, roleRepo: roleRepo}
}