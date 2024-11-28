package presenter

import (
	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavingBookOutput struct {
	contracts.AuditedEntity
	AccountId      primitive.ObjectID ` json:"accountId"`
	Address        contracts.Address  ` json:"address"`
	IdCardNumber   string             ` json:"idCardNumber"`
	Regulations    []Regulation       ` json:"regulations"`
	Balance        float64            `bson:"Balance" json:"balance"`
	PendingBalance float64            `bson:"PendingBalance" json:"pendingBalance"`
	Status         string             ` json:"status"`
}
