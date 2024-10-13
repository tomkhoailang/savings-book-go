package domain

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	AggregateRoot     `bson:",inline" json:",inline"`
	SavingBookId      primitive.ObjectID `bson:"SavingBookId" json:"savingBookId"`
	UserId            primitive.ObjectID `bson:"UserId" json:"userId"`
	Message           string             `bson:"Message" json:"message"`
	IsRead            bool               `bson:"IsRead" json:"isRead"`
	Status            string             `bson:"Status" json:"status"`
	TransactionTicketId primitive.ObjectID `bson:"TransactionTicketId" json:"transactionTicketId"`
}
