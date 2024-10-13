package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionTicket struct {
	AggregateRoot     `bson:",inline" json:",inline"`
	SavingBookId      primitive.ObjectID `bson:"SavingBookId" json:"savingBookId"`
	TransactionDate   time.Time          `bson:"TransactionDate" json:"transactionDate"`
	TransactionType   string             `bson:"TransactionType" json:"transactionType"`
	TransactionMethod string             `bson:"TransactionMethod" json:"transactionMethod"`
	Status            string             `bson:"Status" json:"status"`
	PaymentId         string             `bson:"PaymentId" json:"paymentId"`
	Amount            float64            `bson:"Amount" json:"amount"`
	Email             string             `bson:"Email" json:"email"`
}


