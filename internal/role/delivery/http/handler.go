package http

import (
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"SavingBooks/internal/role/presenter"
	"SavingBooks/utils"
	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	roleUC role.UseCase
}

func (rh *roleHandler) CreateRole() gin.HandlerFunc {
	return utils.HandleCreateRequest[presenter.RoleInput, presenter.RoleOutput, domain.Role](rh.roleUC.CreateRole)
}
func (rh *roleHandler) UpdateRole() gin.HandlerFunc {
	return utils.HandleUpdateRequest[presenter.RoleInput, presenter.RoleOutput, domain.Role](rh.roleUC.UpdateRole)
}

func (rh *roleHandler) DeleteManyRoles() gin.HandlerFunc {
	return utils.HandleDeleteManyRequest[domain.Role](rh.roleUC.DeleteManyRoles)
}

func (rh *roleHandler) GetListRoles() gin.HandlerFunc {
	return utils.HandleGetListRequest[domain.Role](rh.roleUC.GetListRoles)
}

func NewRoleHandler(roleUC role.UseCase) role.Handler {
	return &roleHandler{roleUC: roleUC}
}

