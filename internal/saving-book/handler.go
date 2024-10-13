package saving_book

import "github.com/gin-gonic/gin"

type Handler interface {
	CreateSavingBook() gin.HandlerFunc
	GetListSavingBook() gin.HandlerFunc
	CreateSavingBookGuest() gin.HandlerFunc
	ConfirmPayment() gin.HandlerFunc
	WithDrawOnline() gin.HandlerFunc

}
