package http

import (
	"SavingBooks/internal/domain"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type transactionTicketHandler struct {
	ticketUC transaction_ticket.UseCase
}

func (t *transactionTicketHandler) GetListTicket() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.TransactionTicket](t.ticketUC.GetListTransactionTicket)
}

func NewTransactionTicketHandler(ticketUC transaction_ticket.UseCase) transaction_ticket.Handler {
	return &transactionTicketHandler{ticketUC: ticketUC}
}


