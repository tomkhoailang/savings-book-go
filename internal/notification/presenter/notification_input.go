package presenter

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationInput struct {
	TransactionTicketId primitive.ObjectID ` json:"transactionTicketId" validate:"required"`
	UserId              primitive.ObjectID ` json:"userId" validate:"required"`
	Message             string             ` json:"message" validate:"required"`
	Status              string             ` json:"status" validate:"required"`
}
