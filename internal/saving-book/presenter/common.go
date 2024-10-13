package presenter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Regulation struct {
	RegulationIdRef  primitive.ObjectID ` json:"regulationIdRef"`
	ApplyDate        time.Time          ` json:"applyDate"`
	Name             string             ` json:"name"`
	TermInMonth      int                ` json:"termInMonth"`
	InterestRate     float64            ` json:"interestRate"`
	MinWithDrawValue float64            ` json:"minWithDrawValue"`
	MinWithDrawDay   int                ` json:"minWithDrawDay"`
}