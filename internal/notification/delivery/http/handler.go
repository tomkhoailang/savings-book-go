package http

import (
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type notificationHandler struct {
	notificationUC notification.UseCase
}

func (n *notificationHandler) GetUserNotifications() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.Notification](n.notificationUC.GetUserNotifications)
}

func NewNotificationHandler(notificationUC notification.UseCase) notification.Handler {
	return &notificationHandler{notificationUC: notificationUC}
}