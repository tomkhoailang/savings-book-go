package http

import (
	"errors"
	"net/http"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type authHandler struct {
	useCase auth.UseCase
}

func (ah *authHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		changePasswordReq := presenter.ChangePasswordReq{}
		err := utils.ReadRequest(c, &changePasswordReq)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "failed to find new password in request"})
			return
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "userId not found"})
			return
		}
		err = ah.useCase.ChangePassword(c, userId, changePasswordReq.OldPassword, changePasswordReq.NewPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Change password successfully"})
		return
	}
}

func (ah *authHandler) ConfirmResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
			resetPasswordConfirm := presenter.ResetPasswordConfirm{}
			err := utils.ReadRequest(c, &resetPasswordConfirm)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "password failed"})
				return
			}
			token := c.Query("token")
			if token == "" {
				c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
				return
			}
			err = ah.useCase.ConfirmResetPassword(c, token, resetPasswordConfirm.Password)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Reset password successfully"})
			return
		}

}

func (ah *authHandler) GenerateResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &presenter.ResetPasswordRequest{}

		err := utils.ReadRequest(c, input)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId not found"})
		}
		err = ah.useCase.GenerateResetPassword(c, input.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,gin.H{"message": "Reset password is sent to your email"} )

	}
}

func (ah *authHandler) LogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId not found"})
		}
		err = ah.useCase.Logout(c, userId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authenticated"})
		}
		c.Status(http.StatusNoContent)
	}
}

func (ah *authHandler) RenewAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &presenter.RenewTokenReq{}

		if err := utils.ReadRequest(c, req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId or refreshToken are missing"})
			return
		}

		accessToken, err := ah.useCase.RenewAccessToken(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
		return
	}
}

func (ah *authHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &presenter.SignUpInput{}

		if err := utils.ReadRequest(c, input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := ah.useCase.SignUp(c.Request.Context(), *input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, presenter.SignUpResponse{Id: user.Id.Hex(), Username: user.Username})
		return
	}
}
func (ah *authHandler) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &presenter.LoginInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokens, err := ah.useCase.SignIn(c.Request.Context(), *input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, *tokens)
		return
	}

}

func NewAuthHandler(useCase auth.UseCase) auth.Handler {
	return &authHandler{useCase: useCase}
}
