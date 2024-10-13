package transaction_ticket

import (
	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionTicketRepository interface {
	domain.GenericRepository[domain.TransactionTicket]
	GetMongoClient() *mongo.Client
}
