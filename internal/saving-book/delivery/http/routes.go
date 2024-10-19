package http

import (
	"SavingBooks/internal/auth/middleware"
	saving_book "SavingBooks/internal/saving-book"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup,s saving_book.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation())
	authGroup.POST("",s.CreateSavingBookOnline())
	authGroup.GET("",s.GetListSavingBook())
	authGroup.POST("/confirm-payment",s.ConfirmPayment())
	authGroup.POST("/:id/withdraw-online",s.WithDrawOnline())
	authGroup.POST("/:id/deposit-online",s.DepositOnline())

}