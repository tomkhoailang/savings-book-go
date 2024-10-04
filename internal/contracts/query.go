package contracts

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	Skip    int    `json:"skip" validate:"min=0,max=100"`
	Max     int    `json:"max" validate:"min=25,max=100"`
	Keyword string `json:"keyword"`
	Sort    string `json:"sort"`
}

func(query *Query) QueryBuilder() (bson.M, *options.FindOptions) {
	filter := bson.M{}
	if query.Keyword != "" {
		filter["keyword"] = bson.M{"$regex": query.Keyword, "$options": "i"}
	}

	options := options.Find()
	options.SetSkip(int64(query.Skip))
	options.SetLimit(int64(query.Max))

	if query.Sort != "" {
		options.SetSort(bson.D{{Key: query.Sort, Value: 1}})
	}
	return filter, options

}