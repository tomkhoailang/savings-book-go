package presenter

import (
	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavingBookOutput struct {
	contracts.AuditedEntity
	AccountId           primitive.ObjectID ` json:"accountId"`
	Address             contracts.Address  ` json:"address"`
	IdCardNumber        string             ` json:"idCardNumber"`
	Regulations         []Regulation       ` json:"regulations"`
	Status              string             ` json:"status"`
}
