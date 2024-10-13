package repository

import (
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionTicketRepository struct {
	domain.GenericRepository[domain.TransactionTicket]
	client *mongo.Client
}

func (t *transactionTicketRepository) GetMongoClient() *mongo.Client {
	return  t.client
}

func NewTransactionTicketRepository(client *mongo.Client, db *mongo.Database, collectionName string) transaction_ticket.TransactionTicketRepository {
	baseRepo := contracts.NewBaseRepository[domain.TransactionTicket](db, collectionName)
	return &transactionTicketRepository{
		client: client,
		GenericRepository: baseRepo}
}