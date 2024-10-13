package repository

import (
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"go.mongodb.org/mongo-driver/mongo"
)

type notificationRepository struct {
	domain.GenericRepository[domain.Notification]
}

func NewNotificationRepository(db *mongo.Database, collectionName string) notification.NotificationRepository {
	baseRepo := contracts.NewBaseRepository[domain.Notification](db, collectionName)
	return &notificationRepository{GenericRepository: baseRepo}
}