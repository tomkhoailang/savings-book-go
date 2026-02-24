package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonthlySavingInterest struct {
	AggregateRoot     `bson:",inline" json:",inline"`
	SavingBookId      primitive.ObjectID `bson:"SavingBookId" json:"savingBookId"`
	Amount float64 `bson:"Amount" json:"amount"`
	InterestRate float64 `bson:"InterestRate" json:"interestRate"`
}
