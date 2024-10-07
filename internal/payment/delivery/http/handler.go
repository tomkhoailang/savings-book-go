package http

import (
	"SavingBooks/internal/contracts/paypal"
	"SavingBooks/internal/payment"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentHandler struct {
	paymentUC payment.PaymentUseCase
}

func (p *paymentHandler) TestCaptureOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		res,err := p.paymentUC.CaptureOrder(c,"6LH967419D845083S")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": res})
	}
}

func (p *paymentHandler) TestCreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		payload := &paypal.InitOrderRequest{
			SavingBookId: primitive.NewObjectID().String(),
			Amount:       "150",
		}
		res,err := p.paymentUC.CreateOrder(c,payload)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": res})
	}
}

func (p *paymentHandler) TestSendPayout() gin.HandlerFunc {
	return func(c *gin.Context) {

		payload := &paypal.PayoutRequest{
		}
		res,err := p.paymentUC.SendPayout(c,payload)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": res})
	}
}

//func (p *paymentHandler) TestToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token,err := p.paymentUC.GetPaypalToken()
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"token": token})
//		return
//	}
//
//}

func NewPaymentHandler(paymentUseCase payment.PaymentUseCase) payment.Handler {
	return &paymentHandler{paymentUC: paymentUseCase}
}


