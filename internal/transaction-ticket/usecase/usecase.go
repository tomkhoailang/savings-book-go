package usecase

import (
	"context"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	saving_book "SavingBooks/internal/saving-book"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
)

type transactionTicketUseCase struct {
	ticketRepo     transaction_ticket.TransactionTicketRepository
	savingBookRepo saving_book.SavingBookRepository
}


func (t *transactionTicketUseCase) GetListTransactionTicket(ctx context.Context, query *contracts.Query, auth *presenter.AuthData) (*contracts.QueryResult[domain.TransactionTicket], error) {
	var ticketInterfaces interface{}
	var err error

	if _, ok := auth.Roles["Admin"]; ok {
		ticketInterfaces, err = t.ticketRepo.GetList(ctx, query)
	} else {
		ticketInterfaces, err = t.ticketRepo.GetListAuth(ctx, query, auth.UserId)
	}

	if err != nil {
		return nil, err
	}

	savingBooks := ticketInterfaces.(*contracts.QueryResult[domain.TransactionTicket])
	return savingBooks, nil
}

func (t *transactionTicketUseCase) GetListTransactionTicketOfSavingBook(ctx context.Context, query *contracts.Query, userId , savingBookId string) (*contracts.QueryResult[domain.TransactionTicket], error) {
	var ticketInterfaces interface{}
	var err error
	ticketInterfaces, err = t.ticketRepo.GetListAuthOnReference(ctx, query, userId,"SavingBookId", savingBookId)
	if err != nil {
		return nil, err
	}

	tickets := ticketInterfaces.(*contracts.QueryResult[domain.TransactionTicket])
	return tickets, nil
}


func NewTransactionTicketUseCase( ticketRepo transaction_ticket.TransactionTicketRepository, savingBookRepo saving_book.SavingBookRepository) transaction_ticket.UseCase {
	return &transactionTicketUseCase{ ticketRepo: ticketRepo, savingBookRepo: savingBookRepo}
}
