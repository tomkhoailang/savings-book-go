package domain

type Address struct {
	Country string `bson:"Country" json:"country,omitempty"`
	City    string `bson:"City" json:"city,omitempty"`
	Street  string `bson:"Street" json:"street,omitempty"`
	Zipcode string `bson:"Zipcode" json:"zipcode,omitempty"`
}

