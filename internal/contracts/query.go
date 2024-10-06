package contracts

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	Skip    int    `json:"skip" `
	Max     int    `json:"max" validate:"max=100"`
	Keyword string `json:"keyword"`
	Sort    string `json:"sort"`
}

func(query *Query) QueryBuilder() (bson.M, *options.FindOptions) {
	filter := bson.M{"IsDeleted": false}
	if query.Keyword != "" {
		filter["keyword"] = bson.M{"$regex": query.Keyword, "$options": "i"}
	}
	options := options.Find()
	if query.Max == 0 {
		query.Max = 25
	}
	options.SetSkip(int64(query.Skip))
	options.SetLimit(int64(query.Max))

	if query.Sort != "" {
		options.SetSort(bson.D{{Key: query.Sort, Value: 1}})
	}
	return filter, options

}