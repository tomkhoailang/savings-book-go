package http

import (
	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handler, mw *middleware.MiddleWareManager) {
	authGroup.POST("/register", h.SignUp())
	authGroup.POST("/login", h.SignIn())
	authGroup.POST("/renew-access", h.RenewAccessToken())
	authGroup.POST("/logout", mw.JWTValidation(), h.LogOut())
}
