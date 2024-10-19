package monthly_saving_interest

import (
	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	domain.GenericRepository[domain.MonthlySavingInterest]
	GetMongoClient() *mongo.Client
}
