package contracts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditedEntity struct {
	Id               primitive.ObjectID ` json:"id"`
	ConcurrencyStamp string             ` json:"concurrencyStamp"`
	CreationTime     time.Time          ` json:"creationTime"`
	CreatorId        primitive.ObjectID ` json:"creatorId"`
	CreatorName     string             ` json:"CreatorName"`


	LastModificationTime time.Time          ` json:"lastModificationTime"`
	LastModifierId       primitive.ObjectID ` json:"lastModifierId"`
	LastModifierName     string             ` json:"lastModifierName"`

	IsDeleted    bool               ` json:"isDeleted"`
	DeletionTime time.Time          ` json:"deletionTime"`
	DeleterId    primitive.ObjectID ` json:"deleterId"`
	DeleterName    string ` json:"deleterName"`

	IsActive bool   `json:"isActive"`
}
