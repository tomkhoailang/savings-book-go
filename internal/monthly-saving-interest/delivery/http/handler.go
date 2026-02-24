package http

import (
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type monthlyInterestHandler struct {
	monthlyInterestUC monthly_saving_interest.UseCase
}

func (t *monthlyInterestHandler) GetListMonthlyInterest() gin.HandlerFunc {
	return utils.HandleGetListRequestAuth[domain.MonthlySavingInterest](t.monthlyInterestUC.GetListMonthlyInterest)
}

func NewMonthlyInterestHandler(monthlyInterestUC monthly_saving_interest.UseCase) monthly_saving_interest.Handler {
	return &monthlyInterestHandler{monthlyInterestUC: monthlyInterestUC}
}

