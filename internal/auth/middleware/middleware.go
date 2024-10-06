package middleware

import (
	"errors"
	"net/http"
	"strings"

	"SavingBooks/internal/auth"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type MiddleWareManager struct {
	authUC auth.UseCase
}

func (mw *MiddleWareManager) JWTValidation() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if headerParts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		result, err := mw.authUC.ParseAccessToken(headerParts[1])
		if err != nil {
			if errors.Is(err, auth.ErrInvalidAccessToken) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return

		}
		c.Set("userId",result.UserId)
		c.Set("roles", result.Roles)
		c.Next()
	}
}
func (mw *MiddleWareManager) RoleValidation(reqRoles []string) gin.HandlerFunc{
	return func(c *gin.Context){
		userRoles, err := utils.GetRoles(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		mapReqRoles := utils.SliceToMap[string](reqRoles)
		mapUserRoles := utils.SliceToMap[string](userRoles)

		hasAccess := false
		for userRole := range mapUserRoles{

			if _, ok := mapReqRoles[userRole]; ok {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func NewMiddleWareManager(authUC auth.UseCase) *MiddleWareManager  {
	return &MiddleWareManager{authUC}
}

