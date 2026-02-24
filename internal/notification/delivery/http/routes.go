package http

import (
	"SavingBooks/internal/auth/middleware"
	"SavingBooks/internal/notification"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, n notification.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation())
	authGroup.GET("", n.GetUserNotifications())
	authGroup.PUT("/:id", n.MarkAsReadNotification())
	authGroup.PUT("", n.MarkAsReadAllNotification())
}
