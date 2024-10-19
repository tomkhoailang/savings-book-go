package usecase

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
)

type monthlyUC struct {
	monthlyRepo monthly_saving_interest.Repository
}

func (m *monthlyUC) GetListMonthlyInterestOfSavingBook(ctx context.Context, query *contracts.Query, userId, savingBookId string) (*contracts.QueryResult[domain.MonthlySavingInterest], error) {
	var  monthlyInterfaces interface{}
	var err error
	monthlyInterfaces, err = m.monthlyRepo.GetListAuthOnReference(ctx, query, "","SavingBookId", savingBookId)
	if err != nil {
		return nil, err
	}

	monthly := monthlyInterfaces.(*contracts.QueryResult[domain.MonthlySavingInterest])
	return monthly, nil
}

func NewMonthlyUC(monthlyRepo monthly_saving_interest.Repository) monthly_saving_interest.UseCase {
	return &monthlyUC{monthlyRepo: monthlyRepo}
}