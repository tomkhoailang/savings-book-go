package presenter

import (
	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonthlyInterestOutput struct {
	contracts.AuditedEntity ` json:",inline"`
	SavingBookId            primitive.ObjectID ` json:"savingBookId"`
	Amount                  float64            ` json:"amount"`
	InterestRate            float64            ` json:"interestRate"`
}
