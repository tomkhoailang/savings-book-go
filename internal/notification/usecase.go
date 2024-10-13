package notification

import (
	"context"

	"SavingBooks/internal/notification/presenter"
)

type UseCase interface {
	SendNotification(ctx context.Context, input *presenter.NotificationInput) error
}
