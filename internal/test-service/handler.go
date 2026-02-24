package test_service

import "github.com/gin-gonic/gin"

type Handler interface {
	TestKafkaProducer() gin.HandlerFunc
}
