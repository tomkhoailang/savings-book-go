package domain

import (
	"context"
)

type GenericRepository[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	GetMany(ctx context.Context, ids []string) (*[]T, error)
	GetByField(ctx context.Context, field string, value interface{}) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T, id string, fieldsToUpdate []string) (*T, error)
	Delete(ctx context.Context, deleterId string, id string) error
	DeleteMany(ctx context.Context, deleterId string, ids []string) error
	GetList(ctx context.Context, query interface{}) (interface{}, error)
	GetListAuth(ctx context.Context, query interface{}, currentUserId string) (interface{}, error)
	CountAll(ctx context.Context) (int, error)
	GetCollection() interface{}
}
