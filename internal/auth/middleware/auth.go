package middleware

import (
	"errors"
	"net/http"
	"strings"

	"SavingBooks/internal/auth"
	"github.com/gin-gonic/gin"
)
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
		result, err := mw.authUC.ParseAccessToken(c.Request.Context(), headerParts[1])
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, auth.ErrInvalidAccessToken) {
				status = http.StatusUnauthorized
			}
			c.AbortWithStatus(status)
			return
		}
		c.Set("userId",result.UserId)
		c.Next()
	}
}
