package domain

import (
	"context"
)

type GenericRepository[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, id string, entity *T) error
	Delete(ctx context.Context, id string, deleterId string) error
	DeleteMany(ctx context.Context, ids []string, deleterId string) error
	GetList(ctx context.Context, query interface{}) (interface{}, error)
}
