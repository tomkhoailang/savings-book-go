package notification

import "github.com/gin-gonic/gin"

type Handler interface {
	GetUserNotifications() gin.HandlerFunc
}