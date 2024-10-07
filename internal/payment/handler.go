package payment

import "github.com/gin-gonic/gin"

type Handler interface {
	//TestToken() gin.HandlerFunc
	TestSendPayout() gin.HandlerFunc
	TestCreateOrder() gin.HandlerFunc
	TestCaptureOrder() gin.HandlerFunc
}
