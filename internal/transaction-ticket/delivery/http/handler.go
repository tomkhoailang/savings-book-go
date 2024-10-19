package http

import (
	"net/http"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
	"SavingBooks/internal/transaction-ticket/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type transactionTicketHandler struct {
	ticketUC transaction_ticket.UseCase
}

func (t *transactionTicketHandler) GetListTicket() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.TransactionTicket](t.ticketUC.GetListTransactionTicket)
}
func (t *transactionTicketHandler) GetTicketsOfSavingBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var query contracts.Query
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		savingBookId := c.Param("id")
		if savingBookId == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"err":"id can not be empty"})
			return
		}

		var output contracts.QueryResult[presenter.TransactionTicketOutput]
		res, err := t.ticketUC.GetListTransactionTicketOfSavingBook(c, &query, userId,savingBookId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = copier.Copy(&output, res)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, output)




	}
}
func NewTransactionTicketHandler(ticketUC transaction_ticket.UseCase) transaction_ticket.Handler {
	return &transactionTicketHandler{ticketUC: ticketUC}
}


