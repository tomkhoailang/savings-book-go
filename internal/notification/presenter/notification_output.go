package presenter

import (
	"time"

	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationOutput struct {
	contracts.AuditedEntity     `json:",inline"`
	SavingBookId      primitive.ObjectID `bson:"SavingBookId" json:"savingBookId"`
	UserId            primitive.ObjectID `bson:"UserId" json:"userId"`
	Message           string             `bson:"Message" json:"message"`
	IsRead            bool               `bson:"IsRead" json:"isRead"`
	Status            string             `bson:"Status" json:"status"`
	TransactionDate   time.Time          `bson:"TransactionDate" json:"transactionDate"`
	TransactionType   string             `bson:"TransactionType" json:"transactionType"`
	TransactionMethod string             `bson:"TransactionMethod" json:"transactionMethod"`
	Amount            float64            `bson:"Amount" json:"amount"`
}