package http

import (
	"SavingBooks/internal/auth/middleware"
	"SavingBooks/internal/role"
	"github.com/gin-gonic/gin"
)


func MapAuthRoutes(authGroup *gin.RouterGroup, r role.Handler, mw *middleware.MiddleWareManager) {
	adminOnly := mw.RoleValidation([]string {"Admin"})
	authGroup.Use(mw.JWTValidation())
	authGroup.POST("",adminOnly,r.CreateRole())
	authGroup.PUT("/:id",adminOnly, r.UpdateRole())
	authGroup.DELETE("",adminOnly, r.DeleteManyRoles())
	authGroup.GET("",r.GetListRoles())
}