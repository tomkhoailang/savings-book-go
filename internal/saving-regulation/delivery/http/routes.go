package http

import (
	"SavingBooks/internal/auth/middleware"
	saving_regulation "SavingBooks/internal/saving-regulation"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, s saving_regulation.Handler, mw *middleware.MiddleWareManager) {
	adminOnly := mw.RoleValidation([]string {"Admin"})
	authGroup.Use(mw.JWTValidation())
	authGroup.POST("",adminOnly,s.CreateRegulation())
	authGroup.PUT("/:id",adminOnly, s.UpdateRegulation())
	authGroup.DELETE("",adminOnly, s.DeleteManyRegulations())
	authGroup.GET("",s.GetListRegulations())
}
