package usecase

import (
	"context"

	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"SavingBooks/internal/notification/presenter"
	"SavingBooks/internal/services/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type notificationUseCase struct {
	notificationRepo notification.NotificationRepository
	socket *websocket.Hub
}

func (n *notificationUseCase) SendNotification(ctx context.Context, input *presenter.NotificationInput) error {
	notification := &domain.Notification{
		SavingBookId:        primitive.ObjectID{},
		UserId:              input.UserId,
		Message:             input.Message,
		IsRead:              false,
		Status:              input.Status,
		TransactionTicketId: input.TransactionTicketId,
	}
	notification.SetInit()
	err := n.notificationRepo.Create(ctx, notification)
	n.socket.SendOne(websocket.WithDrawStatus,input.UserId.Hex(), notification)
	if err != nil {
		return err
	}
	return nil
}

func NewNotificationUseCase(notificationRepo notification.NotificationRepository, socket *websocket.Hub) notification.UseCase {
	return &notificationUseCase{notificationRepo: notificationRepo,socket: socket}
}
