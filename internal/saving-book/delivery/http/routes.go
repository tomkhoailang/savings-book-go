package http

import (
	"SavingBooks/internal/auth/middleware"
	saving_book "SavingBooks/internal/saving-book"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup,s saving_book.Handler, mw *middleware.MiddleWareManager) {
	adminOnly := mw.RoleValidation([]string {"Admin"})
	authGroup.Use(mw.JWTValidation())
	authGroup.POST("",adminOnly,s.CreateSavingBookGuest())
	authGroup.GET("",adminOnly,s.GetListSavingBook())
	authGroup.POST("/confirm-payment",adminOnly,s.ConfirmPayment())
	authGroup.POST("/:id/withdraw-online",s.WithDrawOnline())

}