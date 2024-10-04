package http

import (
	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handler, mw *middleware.MiddleWareManager) {
	authGroup.GET("/test",mw.JWTValidation(), func(context *gin.Context) {
		context.JSON(200, "hehe")
	},);
}