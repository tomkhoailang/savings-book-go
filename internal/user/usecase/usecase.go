package usecase

import (
	"context"
	"time"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	"SavingBooks/internal/services/redis"
	"SavingBooks/internal/services/redis/redis_key"
	"SavingBooks/internal/user"
	"SavingBooks/internal/user/presenter"
	"github.com/jinzhu/copier"
)

type userUseCase struct {
	userRepo     auth.UserRepository
	cacheService *redis.Cache
	monthlyRepo monthly_saving_interest.Repository
	expireTime time.Duration
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
	user.IsActive = !user.IsActive
	user.RefreshToken = ""
	_, err = u.userRepo.Update(ctx, user, userId, []string{"IsActive", "RefreshToken"})
	if err != nil {
		return err
	}
	err = u.cacheService.SetValueWithExpire(ctx, redis_key.BlockUserId+":"+userId, user.Id, u.expireTime)
	if err != nil {
		return err
	}
	return nil
}

func NewUserUseCase(userRepo auth.UserRepository, monthlyRepo monthly_saving_interest.Repository,  cacheService *redis.Cache, expireTime time.Duration) user.UseCase {
	return &userUseCase{userRepo: userRepo, cacheService: cacheService, expireTime: expireTime, monthlyRepo: monthlyRepo}
}
