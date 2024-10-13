package http

import (
	"SavingBooks/internal/auth/middleware"
	test_service "SavingBooks/internal/test-service"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, t test_service.Handler, mw *middleware.MiddleWareManager) {
	authGroup.GET("/test", t.TestKafkaProducer());
}