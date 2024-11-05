package auth

import "github.com/gin-gonic/gin"

type Handler interface {
	SignUp() gin.HandlerFunc
	SignIn() gin.HandlerFunc
	ChangePassword() gin.HandlerFunc
	GenerateResetPassword() gin.HandlerFunc
	ConfirmResetPassword() gin.HandlerFunc
	RenewAccessToken() gin.HandlerFunc
	LogOut() gin.HandlerFunc
}
