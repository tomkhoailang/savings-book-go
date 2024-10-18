package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionTicket struct {
	AggregateRoot     `bson:",inline" json:",inline"`
	SavingBookId      primitive.ObjectID `bson:"SavingBookId" json:"savingBookId"`
	TransactionDate   time.Time          `bson:"TransactionDate" json:"transactionDate"`
	Status            string             `bson:"Status" json:"status"`
	Email             string             `bson:"Email" json:"email"`

	PaymentLink  string `bson:"PaymentLink" json:"paymentLink"`
	PaymentType      string  `bson:"PaymentType" json:"paymentType"`
	PaymentId     string  `bson:"PaymentId" json:"paymentId"`
	PaymentAmount float64 `bson:"PaymentAmount" json:"paymentAmount"`
}


