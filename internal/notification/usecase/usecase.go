package usecase

import (
	"context"

	presenter2 "SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"SavingBooks/internal/notification/presenter"
	"SavingBooks/internal/services/websocket"
)

type notificationUseCase struct {
	notificationRepo notification.NotificationRepository
	socket *websocket.Hub
}

func (n *notificationUseCase) GetUserNotifications(ctx context.Context, query *contracts.Query, auth *presenter2.AuthData) (*contracts.QueryResult[domain.Notification], error) {
	var notificationInterface interface{}
	var err error
	notificationInterface, err = n.notificationRepo.GetListAuthOnReference(ctx, query, "","UserId",auth.UserId)
	if err != nil {
		return nil, err
	}

	notifications := notificationInterface.(*contracts.QueryResult[domain.Notification])
	return notifications, nil
}

func (n *notificationUseCase) SendNotification(ctx context.Context, input *presenter.NotificationInput) error {
	notification := &domain.Notification{
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
