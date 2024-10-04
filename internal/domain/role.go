package domain

type Role struct {
	AggregateRoot `bson:",inline" json:",inline"`
	Name string `bson:"Name" json:"name"`
	Description string `bson:"Description" json:"description"`
}
