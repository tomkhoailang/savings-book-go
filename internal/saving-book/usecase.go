package saving_book

import (
	"context"

	presenter2 "SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/saving-book/presenter"
	"SavingBooks/internal/services/kafka/event"
)

type UseCase interface {
	CreateSavingBookOnline(ctx context.Context, input *presenter.SavingBookGuestInput, creatorId string)(*domain.SavingBook, error)
	GetListSavingBook(ctx context.Context, query *contracts.Query, auth *presenter2.AuthData) (*contracts.QueryResult[domain.SavingBook], error)
	ConfirmPaymentOnline(ctx context.Context, paymentId ,userId string) error
	WithdrawOnline(ctx context.Context, input *presenter.WithDrawInput, savingBookId, userId string) error
	DepositOnline(ctx context.Context, input *presenter.DepositInput, savingBookId, userId string)  (*domain.TransactionTicket, error)
	HandleWithdraw(ctx context.Context, input *event.WithDrawEvent) error

}
