package websocket

import (
	"SavingBooks/internal/auth/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, hub *Hub, mw *middleware.MiddleWareManager) {
	authGroup.GET("", mw.JWTValidation(), func(context *gin.Context) {
		ServeWs(hub, context)
	})
}
