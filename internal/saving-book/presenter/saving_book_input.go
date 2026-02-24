package presenter

import (
	"SavingBooks/internal/contracts"
)

type SavingBookGuestInput struct {
	Address      contracts.Address ` json:"address" validate:"required"`
	IdCardNumber string            ` json:"idCardNumber" validate:"required,max=12" `
	Term         int            ` json:"term"`
	NewPaymentAmount      float64            ` json:"newPaymentAmount" validate:"required,min=1"`
}

type ConfirmPaymentInput struct {
	PaymentId string            ` json:"paymentId" validate:"required" `

}



