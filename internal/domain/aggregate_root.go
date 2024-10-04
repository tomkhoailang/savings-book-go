package domain

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AggregateRoot struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConcurrencyStamp string             `bson:"ConcurrencyStamp" json:"concurrencyStamp"`
	CreationTime time.Time          `bson:"CreationTime"  json:"creationTime"`
	CreatorId        primitive.ObjectID `bson:"CreatorId,omitempty" json:"creatorId"`

	LastModificationTime time.Time          `bson:"LastModificationTime" json:"lastModificationTime"`
	LastModifierId       primitive.ObjectID `bson:"LastModifierId,omitempty" json:"lastModifierId"`

	IsDeleted    bool               `bson:"IsDeleted" json:"isDeleted"`
	DeletionTime time.Time          `bson:"DeletionTime" json:"deletionTime"`
	DeleterId    primitive.ObjectID `bson:"DeleterId,omitempty" json:"deleterId"`

	IsActive bool   `bson:"IsActive" json:"isActive"`
	Keyword  string `bson:"Keyword" json:"keyword"`
}
func (a *AggregateRoot) SetInit() {
	a.Id = primitive.NewObjectID()
	a.CreationTime = time.Now()
	a.IsActive = true
	a.IsDeleted = false
	a.ConcurrencyStamp = uuid.New().String()
}
func (a *AggregateRoot) SetCreate(creatorId string) {
	a.Id = primitive.NewObjectID()
	a.CreationTime = time.Now()
	a.CreatorId, _ = primitive.ObjectIDFromHex(creatorId)
	a.IsActive = true
	a.IsDeleted = false
	a.ConcurrencyStamp = uuid.New().String()
}

func (a *AggregateRoot) SetUpdate(lastModifierId string) {
	a.LastModifierId,_ = primitive.ObjectIDFromHex(lastModifierId)
	a.LastModificationTime = time.Now()
	a.ConcurrencyStamp = uuid.New().String()
}

func (a *AggregateRoot) SetDelete(deleterId primitive.ObjectID) {
	a.IsDeleted = true
	a.DeletionTime = time.Now()
	a.DeleterId = deleterId
	a.ConcurrencyStamp = uuid.New().String()
}

