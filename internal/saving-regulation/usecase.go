package saving_regulation

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/saving-regulation/presenter"
)

type SavingRegulationUseCase interface {
	CreateRegulation(ctx context.Context, input *presenter.SavingRegulationInput, creatorId string) (*domain.SavingRegulation, error)
	UpdateRegulation(ctx context.Context, input *presenter.SavingRegulationInput, lastModifierId, regulationId string) (*domain.SavingRegulation, error)
	DeleteManyRegulations(ctx context.Context, deleterId string, input []string) error
	GetListRegulation(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[domain.SavingRegulation], error)
}