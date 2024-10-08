package repository

import (
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/saving-regulation"
	"go.mongodb.org/mongo-driver/mongo"
)

type savingRepository struct {
	domain.GenericRepository[domain.SavingRegulation]
}

func NewSavingRepository(db *mongo.Database, collectionName string) saving_regulation.SavingRegulationRepository {
	baseRepo := contracts.NewBaseRepository[domain.SavingRegulation](db, collectionName)
	return &savingRepository{GenericRepository: baseRepo}
}
