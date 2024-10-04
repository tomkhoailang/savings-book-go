package presenter

import "SavingBooks/internal/contracts"

type RoleOutput struct {
	contracts.AuditedEntity `json:",inline"`
	Name string `json:"name"`
	Description string `json:"description"`
}
