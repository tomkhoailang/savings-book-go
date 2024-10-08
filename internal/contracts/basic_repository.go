package contracts

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository[T any] struct {
	db             *mongo.Database
	collectionName string
	dbName         string
}


func (r *BaseRepository[T]) Get(ctx context.Context, id string) (*T, error) {
	collection := r.db.Collection(r.collectionName)
	var entity T

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
func (r *BaseRepository[T]) GetByField(ctx context.Context, field string, value interface{}) (*T, error) {
	collection := r.db.Collection(r.collectionName)
	var entity T

	var queryValue interface{}
	switch v := value.(type) {
	case string:
		queryValue = v
	case int:
		queryValue = v
	case int64:
		queryValue = v
	case float64:
		queryValue = v
	case bool:
		queryValue = v
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}

	err := collection.FindOne(ctx, bson.M{field: queryValue}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	collection := r.db.Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, entity)
	return err
}
func (r *BaseRepository[T]) GetCollection() *mongo.Collection {
	collection := r.db.Collection(r.collectionName)
	return collection
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, id string, fieldsToUpdate []string) (*T, error) {
	collection := r.db.Collection(r.collectionName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ObjectID: " + id)
	}
	filter := bson.M{"_id": objectId}
	fieldsToUpdate = append(fieldsToUpdate,"LastModifierId", "LastModificationTime", "Keyword", "IsActive")
	updateFields := bson.M{}

	val := reflect.ValueOf(entity).Elem()
	for _, field := range fieldsToUpdate {
		fieldValue := val.FieldByName(field)
		if fieldValue.IsValid() {
			updateFields[field] = fieldValue.Interface()
		}
	}
	if len(updateFields) == 0 {
		return nil, errors.New("no valid fields to update")
	}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedEntity T

	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, DocumentNotFound
		}
		return nil, err
	}
	return &updatedEntity, err
}
func (r *BaseRepository[T]) Delete(ctx context.Context, deleterId string, id string) error {
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
func (r *BaseRepository[T]) DeleteMany(ctx context.Context, deleterId string, ids []string) error {
	collection := r.db.Collection(r.collectionName)

	var objectIDs []primitive.ObjectID

	// Convert string IDs to ObjectIDs
	for _, idStr := range ids {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return errors.New("invalid ObjectID: " + idStr)
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}, "IsDeleted": false}
	update := bson.M{
		"$set": bson.M{
			"DeletionTime": time.Now(),
			"IsDeleted":    true,
			"DeleterId":   deleterId,
		},
	}
	result, err := collection.UpdateMany(ctx, filter, update)
	if result != nil && result.ModifiedCount == 0 {
		return DocumentNotFound
	}
	return err
}
func (r *BaseRepository[T]) CountAll(ctx context.Context) (int, error) {
	collection := r.db.Collection(r.collectionName)
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return -1, err
	}
	return int(totalCount),nil
}
func (r *BaseRepository[T]) GetList(ctx context.Context, query interface{}) (interface{}, error) {
	collection := r.db.Collection(r.collectionName)

	filter, options := query.(*Query).QueryBuilder()

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

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
func (r *BaseRepository[T]) GetMany(ctx context.Context, ids []string) (*[]T, error) {
	collection := r.db.Collection(r.collectionName)

	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New("invalid ObjectID: " + id)
		}
		objectIds = append(objectIds, objectId)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var entities []T
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return &entities, nil
}

func NewBaseRepository[T any](db *mongo.Database, collectionName string) domain.GenericRepository[T]  {
	return &BaseRepository[T]{db: db, collectionName: collectionName}
}
