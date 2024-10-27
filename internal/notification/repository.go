package notification

import (
	"context"

	"SavingBooks/internal/domain"
)

type NotificationRepository interface {
	domain.GenericRepository[domain.Notification]
	MarkAsReadAllNotification(ctx context.Context, userId string) error
}
