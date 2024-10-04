package http

import (
	"SavingBooks/internal/auth"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handler) {
	authGroup.POST("/register", h.SignUp());
	authGroup.POST("/login", h.SignIn());
}