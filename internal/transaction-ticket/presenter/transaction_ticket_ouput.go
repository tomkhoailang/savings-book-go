package presenter

import (
	"time"

	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionTicket struct {
	contracts.AuditedEntity
	SavingBookId    primitive.ObjectID ` json:"savingBookId"`
	TransactionDate time.Time          ` json:"transactionDate"`
	Status          string             ` json:"status"`
	Email           string             ` json:"email"`

	PaymentLink   string  ` json:"paymentLink"`
	PaymentType   string  ` json:"paymentType"`
	PaymentId     string  ` json:"paymentId"`
	PaymentAmount float64 ` json:"paymentAmount"`
}
