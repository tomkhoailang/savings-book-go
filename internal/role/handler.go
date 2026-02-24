package role

import "github.com/gin-gonic/gin"

type Handler interface {
	CreateRole() gin.HandlerFunc
	UpdateRole() gin.HandlerFunc
	DeleteManyRoles() gin.HandlerFunc
	GetListRoles() gin.HandlerFunc
}
