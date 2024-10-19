package transaction_ticket

import "github.com/gin-gonic/gin"

type Handler interface {
	GetListTicket() gin.HandlerFunc
	GetTicketsOfSavingBook() gin.HandlerFunc
}
