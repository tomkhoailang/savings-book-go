package presenter

type SavingRegulationInput struct {
	MinWithdrawValue float64      `json:"minWithdrawValue" validate:"min=10"`
	SavingTypes      []SavingType `json:"savingTypes"`
	MinWithdrawDay   int          `json:"minWithdrawDay"`
	IsActive         bool         `json:"isActive"`
}

type SavingType struct {
	Name         string  `json:"name"`
	Term         int     `json:"term" validate:"min=1"`
	InterestRate float64 `json:"interestRate"`
}
