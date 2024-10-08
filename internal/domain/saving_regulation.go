package domain

type SavingRegulation struct {
	AggregateRoot `bson:",inline" json:",inline"`
	MinWithdrawValue float64      `bson:"MinWithdrawValue" json:"minWithdrawValue" validate:"min=10"`
	SavingTypes      []SavingType `bson:"SavingTypes" json:"savingTypes"`
	MinWithdrawDay   int          `bson:"MinWithdrawDay" json:"minWithdrawDay"`
	IsActive         bool         `bson:"IsActive" json:"isActive"`
}
type SavingType struct {
	Name         string  `json:"name"`
	Term         int     `json:"term"`
	InterestRate float64 `json:"interestRate"`
}
