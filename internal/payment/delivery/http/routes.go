package http

import (
	"SavingBooks/internal/payment"
	"github.com/gin-gonic/gin"
)


func MapAuthRoutes(authGroup *gin.RouterGroup, h payment.Handler) {
	authGroup.POST("/send-payout", h.TestSendPayout())
	authGroup.POST("/create-order", h.TestCreateOrder())
	authGroup.POST("/capture-order", h.TestCaptureOrder())
}