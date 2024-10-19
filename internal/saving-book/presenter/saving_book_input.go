package presenter

import (
	"SavingBooks/internal/contracts"
)

type SavingBookGuestInput struct {
	Address      contracts.Address ` json:"address" validate:"required"`
	IdCardNumber string            ` json:"idCardNumber" validate:"required,min=12" `
	Term         int            ` json:"term" validate:"required"`
	NewPaymentAmount      float64            ` json:"newPaymentAmount" validate:"required,min=10"`
}

type ConfirmPaymentInput struct {
	PaymentId string            ` json:"paymentId" validate:"required" `

}



