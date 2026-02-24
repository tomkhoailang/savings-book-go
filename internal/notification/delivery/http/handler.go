package http

import (
	"net/http"

	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type notificationHandler struct {
	notificationUC notification.UseCase
}

func (n *notificationHandler) MarkAsReadNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		notificationId := c.Param("id")
		if notificationId == "" {
			c.JSON(http.StatusNoContent, gin.H{})
			return
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{})
			return
		}
		err = n.notificationUC.MarkAsReadNotification(c, userId, notificationId)
		c.JSON(http.StatusNoContent, gin.H{})
		return

	}
}
func (n *notificationHandler) MarkAsReadAllNotification() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{})
			return
		}
		err = n.notificationUC.MarkAsReadAllNotification(c, userId)
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
}

func (n *notificationHandler) GetUserNotifications() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.Notification](n.notificationUC.GetUserNotifications)
}

func NewNotificationHandler(notificationUC notification.UseCase) notification.Handler {
	return &notificationHandler{notificationUC: notificationUC}
}