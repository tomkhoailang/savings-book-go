package usecase

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	saving_regulation "SavingBooks/internal/saving-regulation"
	"SavingBooks/internal/saving-regulation/presenter"
	"SavingBooks/internal/services/redis"
	"SavingBooks/internal/services/redis/redis_key"
	"github.com/jinzhu/copier"
)

type savingRegulationUseCase struct {
	savingRegulationRepo saving_regulation.SavingRegulationRepository
	cacheService *redis.Cache
}

func (s *savingRegulationUseCase) CreateRegulation(ctx context.Context, input *presenter.SavingRegulationInput, creatorId string) (*domain.SavingRegulation, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	regulation := &domain.SavingRegulation{}
	err = copier.Copy(regulation, input)
	if err != nil {
		return nil, err
	}
	regulation.SetCreate(creatorId)
	err = s.savingRegulationRepo.Create(ctx, regulation)

	if err != nil {
		return nil, err
	}
	_ = s.cacheService.RemoveValue(ctx, redis_key.LatestRegulation)

	return regulation, nil
}

func (s *savingRegulationUseCase) UpdateRegulation(ctx context.Context, input *presenter.SavingRegulationInput, lastModifierId, regulationId string) (*domain.SavingRegulation, error) {
	regulation := &domain.SavingRegulation{}
	err := copier.CopyWithOption(regulation, input, copier.Option{IgnoreEmpty: true,DeepCopy: true})
	if err != nil {
		return nil, err
	}
	regulation.SetUpdate(lastModifierId)
	regulation, err = s.savingRegulationRepo.Update(ctx, regulation, regulationId, []string{"MinWithdrawDay", "SavingTypes", "MinWithdrawValue"})
	if err != nil {
		return nil, err
	}
	_ = s.cacheService.RemoveValue(ctx, redis_key.LatestRegulation)
	return regulation, nil
}

func (s *savingRegulationUseCase) DeleteManyRegulations(ctx context.Context, deleterId string, input []string) error {
	err := s.savingRegulationRepo.DeleteMany(ctx, deleterId, input)
	_ = s.cacheService.RemoveValue(ctx, redis_key.LatestRegulation)
	if err != nil {
		return err
	}
	return nil
}

func (s *savingRegulationUseCase) GetListRegulation(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[domain.SavingRegulation], error) {

	regulationsInterface, err := s.savingRegulationRepo.GetList(ctx, query)
	if err != nil {
		return nil, err
	}
	regulations := regulationsInterface.(*contracts.QueryResult[domain.SavingRegulation])
	return regulations, nil

}
func (s *savingRegulationUseCase) GetLatestSavingRegulation(ctx context.Context) (*domain.SavingRegulation, error) {
	return s.cacheService.GetLatestSavingRegulation(ctx)
}

func NewSavingRegulationUseCase(savingRegulationRepo saving_regulation.SavingRegulationRepository, cacheService *redis.Cache) saving_regulation.SavingRegulationUseCase {
	return &savingRegulationUseCase{savingRegulationRepo: savingRegulationRepo, cacheService: cacheService}
}