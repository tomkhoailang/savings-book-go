package saving_regulation

import "github.com/gin-gonic/gin"

type Handler interface {
	CreateRegulation() gin.HandlerFunc
	UpdateRegulation() gin.HandlerFunc
	DeleteManyRegulations() gin.HandlerFunc
	GetListRegulations() gin.HandlerFunc
}
