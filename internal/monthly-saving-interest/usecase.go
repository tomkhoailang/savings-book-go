package monthly_saving_interest

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	GetListMonthlyInterestOfSavingBook(ctx context.Context, query *contracts.Query, userId , savingBookId string) (*contracts.QueryResult[domain.MonthlySavingInterest], error)
	GetListMonthlyInterest(ctx context.Context, query *contracts.Query, auth *presenter.AuthData) (*contracts.QueryResult[domain.MonthlySavingInterest], error)
	GetTotalEarningsOfSavingBooks(ctx context.Context, savingBookIds []string) (map[string]float64, error)
}
