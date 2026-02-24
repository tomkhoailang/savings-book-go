package repository

import (
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	saving_book "SavingBooks/internal/saving-book"
	"go.mongodb.org/mongo-driver/mongo"
)

type savingBookRepository struct {
	domain.GenericRepository[domain.SavingBook]
}



func NewSavingBookRepository(db *mongo.Database, collectionName string) saving_book.SavingBookRepository {
	baseRepo := contracts.NewBaseRepository[domain.SavingBook](db, collectionName).(*contracts.BaseRepository[domain.SavingBook])
	return &savingBookRepository{GenericRepository: baseRepo}
}