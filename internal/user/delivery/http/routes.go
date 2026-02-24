package http

import (
	"SavingBooks/internal/auth/middleware"
	"SavingBooks/internal/user"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, u user.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation(), mw.RoleValidation([]string{"Admin"}))
	authGroup.GET("", u.GetListUser())
	authGroup.PUT("disable/:id", u.DisableUser())
}