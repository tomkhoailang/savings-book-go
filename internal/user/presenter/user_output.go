package presenter

import (

	"SavingBooks/internal/contracts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	contracts.AuditedEntity `bson:",inline" json:",inline"`
	Username string `bson:"Username" json:"username"`
	Email string `bson:"Email" json:"email"`
	TotalEarnings float64 `bson:"TotalEarnings" json:"totalEarnings"`
	RoleIds []primitive.ObjectID `bson:"RoleIds" json:"rolesIds"`

}
