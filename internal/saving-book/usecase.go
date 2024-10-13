package saving_book

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/saving-book/presenter"
	"SavingBooks/internal/services/kafka/event"
)

type UseCase interface {
	CreateSavingBook(ctx context.Context)
	CreateSavingBookOnline(ctx context.Context, input *presenter.SavingBookGuestInput, creatorId string)(*domain.SavingBook, error)
	ConfirmPaymentOnline(ctx context.Context, paymentId string) error
	WithdrawOnline(ctx context.Context, input *presenter.WithDrawInput, savingBookId string) error
	HandleWithdraw(ctx context.Context, input *event.WithDrawEvent) error
	GetListSavingRegulation(ctx context.Context, query *contracts.Query) (*contracts.QueryResult[domain.SavingBook], error)

}
