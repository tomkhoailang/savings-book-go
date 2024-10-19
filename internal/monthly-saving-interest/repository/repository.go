package repository

import (
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	"go.mongodb.org/mongo-driver/mongo"
)

type monthlySavingInterestRepository struct {
	domain.GenericRepository[domain.MonthlySavingInterest]
	client *mongo.Client
}

func (t *monthlySavingInterestRepository) GetMongoClient() *mongo.Client {
	return  t.client
}
func NewMonthlySavingInterestRepository(client *mongo.Client, db *mongo.Database, collectionName string) monthly_saving_interest.Repository {
	baseRepo := contracts.NewBaseRepository[domain.MonthlySavingInterest](db, collectionName)
	return &monthlySavingInterestRepository{GenericRepository: baseRepo, client: client}
}

