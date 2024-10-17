package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavingBook struct {
	AggregateRoot     `bson:",inline" json:",inline"`
	AccountId         primitive.ObjectID `bson:"AccountId" json:"accountId"`
	Address           Address            `bson:"Address" json:"address"`
	IdCardNumber      string             `bson:"IdCardNumber" json:"idCardNumber"`
	Regulations       []Regulation       `bson:"Regulations" json:"regulations"`
	Balance           float64            `bson:"Balance" json:"balance"`
	Status            string             `bson:"Status" json:"status"`
	NextScheduleMonth time.Time          `bson:"NextScheduleMonth" json:"nextScheduleMonth"`

	NewPaymentLink   string  `bson:"NewPaymentLink" json:"newPaymentLink"`
	NewPaymentType   string  `bson:"NewPaymentType" json:"newPaymentType"`
	NewPaymentId     string  `bson:"NewPaymentId" json:"newPaymentId"`
	NewPaymentAmount float64 `bson:"NewPaymentAmount" json:"newPaymentAmount"`
}

type Regulation struct {
	RegulationIdRef  primitive.ObjectID `bson:"RegulationIdRef" json:"regulationIdRef"`
	ApplyDate        time.Time          `bson:"ApplyDate" json:"applyDate"`
	Name             string             `bson:"Name" json:"name"`
	TermInMonth      int                `bson:"TermInMonth" json:"termInMonth"`
	InterestRate     float64            `bson:"InterestRate" json:"interestRate"`
	MinWithDrawValue float64            `bson:"MinWithDrawValue" json:"minWithDrawValue"`
	MinWithDrawDay   int                `bson:"MinWithDrawDay" json:"minWithDrawDay"`
	NoTermInterestRate float64 `bson:"NoTermInterestRate" json:"noTermInterestRate"`
}
