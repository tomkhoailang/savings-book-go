package http

import (
	"SavingBooks/internal/domain"
	"SavingBooks/internal/user"
	"SavingBooks/internal/user/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type userHandler struct {
	userUC user.UseCase
}

func (u *userHandler) GetListUser() gin.HandlerFunc {
	return utils.NewHandleGetListRequest[domain.User, presenter.User](u.userUC.GetListUser)
}

func (u *userHandler) DisableUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		if userId == "" {
			c.JSON(500, gin.H{"err": "id not found"})
			return
		}
		err := u.userUC.DisableUser(c, userId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(404, gin.H{"err": "user not found"})
				return
			}
			c.JSON(500, gin.H{"err": err.Error()})
			return
		}
		c.JSON(200, "ok")
		return
	}
}

func NewUserHandler(userUC user.UseCase) user.Handler {
	return &userHandler{userUC: userUC}
}



