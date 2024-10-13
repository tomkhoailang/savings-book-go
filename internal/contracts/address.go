package contracts

type Address struct {
	Country string ` json:"country,omitempty"`
	City    string ` json:"city,omitempty"`
	Street  string ` json:"street,omitempty"`
	Zipcode string ` json:"zipcode,omitempty"`
}
