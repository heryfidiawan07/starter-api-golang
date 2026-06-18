package handler

import (
	"errors"
	"math"
	"strconv"

	"starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/internal/service"
	"starter-api-golang/pkg/response"
	"starter-api-golang/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	roleUsecase usecase.RoleUsecase
}

func NewRoleHandler(roleUsecase usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase: roleUsecase}
}

func (h *RoleHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")

	filter := repository.RoleFilter{
		Search:  search,
		Page:    page,
		PerPage: perPage,
	}

	roles, total, err := h.roleUsecase.FindAll(filter)
	if err != nil {
		response.InternalServerError(c, "failed to retrieve roles")
		return
	}

	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	response.OKWithMeta(c, "roles retrieved", roles, &response.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPage,
	})
}

func (h *RoleHandler) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid role ID", nil)
		return
	}

	role, err := h.roleUsecase.FindByID(id)
	if err != nil {
		response.NotFound(c, "role not found")
		return
	}

	response.OK(c, "role retrieved", role)
}

func (h *RoleHandler) Store(c *gin.Context) {
	var req usecase.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	role, err := h.roleUsecase.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNameTaken):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "failed to create role")
		}
		return
	}

	response.Created(c, "role created", role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid role ID", nil)
		return
	}

	var req usecase.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	role, err := h.roleUsecase.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNotFound):
			response.NotFound(c, "role not found")
		case errors.Is(err, service.ErrRoleNameTaken):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "failed to update role")
		}
		return
	}

	response.OK(c, "role updated", role)
}

func (h *RoleHandler) Destroy(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid role ID", nil)
		return
	}

	if err := h.roleUsecase.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNotFound):
			response.NotFound(c, "role not found")
		case errors.Is(err, service.ErrRoleInUse):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "failed to delete role")
		}
		return
	}

	response.OK(c, "role deleted", nil)
}
