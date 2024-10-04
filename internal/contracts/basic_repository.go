package contracts

import (
	"context"
	"errors"
	"time"

	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository[T any] struct {
	db             *mongo.Database
	collectionName string
	dbName         string
}

func NewBaseRepository[T any](db *mongo.Database, collectionName string) domain.GenericRepository[T]  {
	return &BaseRepository[T]{db: db, collectionName: collectionName}
}

func (r *BaseRepository[T]) Get(ctx context.Context, id string) (*T, error) {
	collection := r.db.Collection(r.collectionName)
	var entity T
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	collection := r.db.Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, entity)
	return err
}

func (r *BaseRepository[T]) Update(ctx context.Context, id string, entity *T) error {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": id}
	_, err := collection.ReplaceOne(ctx, filter, entity)
	return err
}
func (r *BaseRepository[T]) Delete(ctx context.Context, id string, deleterId string) error {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"DeletionTime": time.Now(),
			"IsDeleted":    true,
			"DeleterId":   deleterId,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *BaseRepository[T]) DeleteMany(ctx context.Context, ids []string, deleterId string) error {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{
		"$set": bson.M{
			"DeletionTime": time.Now(),
			"IsDeleted":    true,
			"DeleterId":   deleterId,
		},
	}
	_, err := collection.UpdateMany(ctx, filter, update)
	return err
}
func (r *BaseRepository[T]) GetList(ctx context.Context, query interface{}) (interface{}, error) {
	collection := r.db.Collection(r.collectionName)

	filter, options := query.(*Query).QueryBuilder()

	totalCount, err := collection.CountDocuments(ctx, filter)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []T
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	queryResult := &QueryResult[T]{
		TotalCount: int(totalCount),
		Items:      entities,
	}

	return queryResult, nil
}
