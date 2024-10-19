package presenter

import (
	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavingBookInput struct {
	AccountId    primitive.ObjectID ` json:"accountId" validate:"required"`
	Address      contracts.Address  ` json:"address" validate:"required"`
	IdCardNumber string             ` json:"idCardNumber" validate:"required"`
	Regulations  string             ` json:"regulationId" validate:"required"`
	NewPaymentAmount      float64            ` json:"newPaymentAmount" validate:"required,min=10"`
}

type SavingBookGuestInput struct {
	Address      contracts.Address ` json:"address" validate:"required"`
	IdCardNumber string            ` json:"idCardNumber" validate:"required,min=12" `
	Term         int            ` json:"term" validate:"required"`
	NewPaymentAmount      float64            ` json:"newPaymentAmount" validate:"required,min=10"`
}

type ConfirmPaymentInput struct {
	PaymentId string            ` json:"paymentId" validate:"required" `

}

