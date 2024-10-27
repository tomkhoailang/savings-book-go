package repository

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type notificationRepository struct {
	domain.GenericRepository[domain.Notification]
}

func (n *notificationRepository) MarkAsReadAllNotification(ctx context.Context, userId string) error {
	collectionInterface := n.GetCollection()

	collection := collectionInterface.(*mongo.Collection)

	nUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.M{
		"UserId": nUserId,
	}
	update := bson.M{
		"$set": bson.M{
			"IsRead": true,
		},
	}
	_,err = collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func NewNotificationRepository(db *mongo.Database, collectionName string) notification.NotificationRepository {
	baseRepo := contracts.NewBaseRepository[domain.Notification](db, collectionName)
	return &notificationRepository{GenericRepository: baseRepo}
}