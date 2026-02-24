package notification

import (
	"context"

	presenter2 "SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification/presenter"
)

type UseCase interface {
	SendNotification(ctx context.Context, input *presenter.NotificationInput) error
	GetUserNotifications(ctx context.Context, query *contracts.Query, auth *presenter2.AuthData) (*contracts.QueryResult[domain.Notification], error)
	MarkAsReadNotification(ctx context.Context, userId, notificationId string) error
	MarkAsReadAllNotification(ctx context.Context, userId string) error
}
