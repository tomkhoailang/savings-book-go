package http

import (
	"net/http"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	useCase auth.UseCase
}

func NewAuthHandler(useCase auth.UseCase) auth.Handler {
	return &authHandler{useCase: useCase}
}


func (ah *authHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &presenter.SignUpInput{}

		if err := utils.ReadRequest(c, input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := ah.useCase.SignUp(c.Request.Context(), *input);
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, presenter.SignUpResponse{Id: user.Id.Hex(), Username: user.Username})
		return
	}
}
func (ah *authHandler) SignIn()  gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &presenter.LoginInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := ah.useCase.SignIn(c.Request.Context(), *input);
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, presenter.LogInResponse{Token: token })
		return
	}

}










