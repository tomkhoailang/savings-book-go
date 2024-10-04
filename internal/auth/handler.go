package auth

import "github.com/gin-gonic/gin"

type Handler interface {
	SignUp() gin.HandlerFunc
	SignIn() gin.HandlerFunc

}
