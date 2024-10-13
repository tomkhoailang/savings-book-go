package notification

import "SavingBooks/internal/domain"

type NotificationRepository interface {
	domain.GenericRepository[domain.Notification]
}
