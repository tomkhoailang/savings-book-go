package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error       string   `json:"error"`
	Message     string   `json:"message"`
	StackTrace  []string `json:"stackTrace,omitempty"`
	Code        int      `json:"code"`
	RequestPath string   `json:"requestPath"`
}

func CustomRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				stackLines := strings.Split(string(stack), "\n")

				var cleanedStack []string
				for _, line := range stackLines {
					if line != "" {
						cleanedStack = append(cleanedStack, strings.TrimSpace(line))
						if len(cleanedStack) == 10 {
							break
						}
					}
				}

				errorResponse := ErrorResponse{
					Error:       "Internal Server Error",
					Message:     fmt.Sprintf("%v", err),
					StackTrace:  cleanedStack,
					Code:        http.StatusInternalServerError,
					RequestPath: c.Request.URL.Path,
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse)
			}
		}()

		c.Next()
	}
}