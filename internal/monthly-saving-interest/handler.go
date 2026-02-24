package monthly_saving_interest

import "github.com/gin-gonic/gin"

type Handler interface {
	GetListMonthlyInterest() gin.HandlerFunc
}
