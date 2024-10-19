package monthly_saving_interest

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	GetListMonthlyInterestOfSavingBook(ctx context.Context, query *contracts.Query, userId , savingBookId string) (*contracts.QueryResult[domain.MonthlySavingInterest], error)
}
