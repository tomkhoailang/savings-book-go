package http

import (
	"SavingBooks/internal/auth/middleware"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, t monthly_saving_interest.Handler, mw *middleware.MiddleWareManager) {
	authGroup.Use(mw.JWTValidation())
	authGroup.GET("", t.GetListMonthlyInterest())
}
