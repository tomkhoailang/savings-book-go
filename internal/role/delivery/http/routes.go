package http

import (
	"net/http"

	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"SavingBooks/internal/role/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type roleHandler struct {
	roleUC role.UseCase
}


func (rh *roleHandler) CreateRole() gin.HandlerFunc {
	return utils.HandleCreateRequest[presenter.RoleInput, presenter.RoleOutput, domain.Role](rh.roleUC.CreateRole)
}

func (rh *roleHandler) UpdateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input presenter.RoleInput
		err := utils.ReadRequest(c, &input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err":err.Error()})
		}
		userId, err := utils.GetUserId(c)
		if err != nil {
			return
		}
		var output presenter.RoleOutput
		role, err := rh.roleUC.UpdateRole(c.Request.Context(),&input, userId)

		err = copier.Copy(output, &role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err":err.Error()})
		}
		c.JSON(http.StatusOK, &output)
		return



	}
}

func (rh *roleHandler) DeleteManyRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "UpdateRole not implemented",
		})
	}
}

func (rh *roleHandler) GetListRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "UpdateRole not implemented",
		})
	}
}

func NewRoleHandler(roleUC role.UseCase) role.Handler {
	return &roleHandler{roleUC: roleUC}
}

//func (rh *roleHandler) CreateRole() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		input := &presenter.RoleInput{}
//		if err := utils.ReadRequest(c, input); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		userId, err := utils.GetUserId(c)
//		if err != nil {
//			return
//		}
//
//		role, err := rh.roleUC.CreateRole(c.Request.Context(), input, userId)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		output := &presenter.RoleOutput{}
//		err = copier.Copy(output, role)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//
//		c.JSON(http.StatusCreated, &output)
//		return
//	}
//}