package handler

import (
	"errors"
	"math"
	"strconv"

	"starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/internal/middleware"
	"starter-api-golang/internal/service"
	"starter-api-golang/pkg/response"
	"starter-api-golang/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")

	filter := repository.UserFilter{
		Search:  search,
		Page:    page,
		PerPage: perPage,
	}

	users, total, err := h.userUsecase.FindAll(filter)
	if err != nil {
		response.InternalServerError(c, "failed to retrieve users")
		return
	}

	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	response.OKWithMeta(c, "users retrieved", users, &response.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPage,
	})
}

func (h *UserHandler) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", nil)
		return
	}

	user, err := h.userUsecase.FindByID(id)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.OK(c, "user retrieved", user)
}

func (h *UserHandler) Store(c *gin.Context) {
	var req usecase.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	user, err := h.userUsecase.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserEmailTaken):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, service.ErrUserUsernameTaken):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "failed to create user")
		}
		return
	}

	response.Created(c, "user created", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", nil)
		return
	}

	var req usecase.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	user, err := h.userUsecase.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFoundSvc):
			response.NotFound(c, "user not found")
		default:
			response.InternalServerError(c, "failed to update user")
		}
		return
	}

	response.OK(c, "user updated", user)
}

func (h *UserHandler) Destroy(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", nil)
		return
	}

	if err := h.userUsecase.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFoundSvc):
			response.NotFound(c, "user not found")
		default:
			response.InternalServerError(c, "failed to delete user")
		}
		return
	}

	response.OK(c, "user deleted", nil)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req usecase.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userUsecase.UpdateProfile(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserUsernameTaken):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "failed to update profile")
		}
		return
	}

	response.OK(c, "profile updated", user)
}

func (h *UserHandler) UpdatePhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "photo file is required", nil)
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userUsecase.UpdatePhoto(userID, file)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "photo updated", user)
}
