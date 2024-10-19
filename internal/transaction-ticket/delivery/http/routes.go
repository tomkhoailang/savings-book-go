package http

import (
	"SavingBooks/internal/auth/middleware"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, t transaction_ticket.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation(), mw.RoleValidation([]string{"Admin"}))
	authGroup.GET("", t.GetListTicket())
}
