package saving_book

import "github.com/gin-gonic/gin"

type Handler interface {
	GetListSavingBook() gin.HandlerFunc
	GetTicketsOfSavingBook() gin.HandlerFunc
	GetMonthlyInterestOfSavingBook() gin.HandlerFunc
	CreateSavingBookOnline() gin.HandlerFunc
	ConfirmPayment() gin.HandlerFunc
	WithDrawOnline() gin.HandlerFunc
	DepositOnline() gin.HandlerFunc

}
