package presenter

import "SavingBooks/internal/contracts"

type SavingRegulationOutput struct {
	contracts.AuditedEntity
	MinWithdrawValue float64      `json:"minWithdrawValue"`
	SavingTypes      []SavingType `json:"savingTypes"`
	MinWithdrawDay   int          `json:"minWithdrawDay"`
}

