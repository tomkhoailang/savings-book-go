package user

import "github.com/gin-gonic/gin"


type Handler interface {
	GetListUser() gin.HandlerFunc
	DisableUser() gin.HandlerFunc
}
