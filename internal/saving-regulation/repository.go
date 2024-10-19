package saving_regulation

import (
	"context"

	"SavingBooks/internal/domain"
)

type SavingRegulationRepository interface {
	domain.GenericRepository[domain.SavingRegulation]
	GetLatestSavingRegulation(ctx context.Context) (*domain.SavingRegulation, error)
}
