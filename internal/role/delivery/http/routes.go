package http

import (
	"SavingBooks/internal/auth/middleware"
	"SavingBooks/internal/role"
	"github.com/gin-gonic/gin"
)


func MapAuthRoutes(authGroup *gin.RouterGroup, r role.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation(), mw.RoleValidation([]string {"Admin"}))
	authGroup.POST("",r.CreateRole())
	authGroup.PUT("/:id", r.UpdateRole())
	authGroup.DELETE("", r.DeleteManyRoles())
	authGroup.GET("",r.GetListRoles())
}