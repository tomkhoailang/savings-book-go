package transaction_ticket

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
)

type UseCase interface {
	GetListTransactionTicket(ctx context.Context, query *contracts.Query, auth *presenter.AuthData) (*contracts.QueryResult[domain.TransactionTicket], error)
	GetTransactionTicket(ctx context.Context, transactionTicketId string, auth *presenter.AuthData) (*domain.TransactionTicket, error)
	GetListTransactionTicketOfSavingBook(ctx context.Context, query *contracts.Query, auth *presenter.AuthData , savingBookId string) (*contracts.QueryResult[domain.TransactionTicket], error)

}
