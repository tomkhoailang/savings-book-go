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
	Balance             float64            ` json:"balance"`
	Status              string             ` json:"status"`
	NewPaymentLink         string             `json:"newPaymentLink"`
	NewPaymentType         string             `json:"newPaymentType"`
	NewPaymentId     string             `json:"newPaymentId"`
	NewPaymentAmount float64            ` json:"newPaymentAmount"`
}
