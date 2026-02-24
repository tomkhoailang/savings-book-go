package repository

import (
	"context"
	"errors"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/saving-regulation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type savingRepository struct {
	domain.GenericRepository[domain.SavingRegulation]
}

func (s *savingRepository) GetLatestSavingRegulation(ctx context.Context) (*domain.SavingRegulation, error) {
	collectionInterface := s.GetCollection()
	collection := collectionInterface.(*mongo.Collection)

	filter := bson.M{"IsActive": true, "IsDeleted": false}
	reg := &domain.SavingRegulation{}
	opts := options.FindOne().SetSort(bson.D{{"CreationTime", -1}})

	err := collection.FindOne(ctx, filter, opts).Decode(&reg)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(saving_regulation.SavingRegulationNotFoundError)
		}
		return nil, err
	}
	return reg, nil

}

func NewSavingRepository(db *mongo.Database, collectionName string) saving_regulation.SavingRegulationRepository {
	baseRepo := contracts.NewBaseRepository[domain.SavingRegulation](db, collectionName)
	return &savingRepository{GenericRepository: baseRepo}
}
