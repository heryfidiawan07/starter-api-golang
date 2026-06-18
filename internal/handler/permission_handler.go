package handler

import (
	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionHandler struct {
	permissionUsecase usecase.PermissionUsecase
}

func NewPermissionHandler(permissionUsecase usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{permissionUsecase: permissionUsecase}
}

func (h *PermissionHandler) Index(c *gin.Context) {
	permissions, err := h.permissionUsecase.FindAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve permissions")
		return
	}
	response.OK(c, "permissions retrieved", permissions)
}

func (h *PermissionHandler) Tree(c *gin.Context) {
	tree, err := h.permissionUsecase.FindTree()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve permission tree")
		return
	}
	response.OK(c, "permission tree retrieved", tree)
}

func (h *PermissionHandler) ByRole(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("role_id"))
	if err != nil {
		response.BadRequest(c, "invalid role ID", nil)
		return
	}

	permissions, err := h.permissionUsecase.FindByRoleID(roleID)
	if err != nil {
		response.InternalServerError(c, "failed to retrieve permissions")
		return
	}

	response.OK(c, "permissions retrieved", permissions)
}
