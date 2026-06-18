package usecase

import (
	"mime/multipart"
	"starter-api-golang/internal/domain/entity"
	"starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string     `json:"username" validate:"required,min=3,max=50"`
	Name     string     `json:"name" validate:"required,min=2,max=100"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=8"`
	RoleID   *uuid.UUID `json:"role_id"`
	IsActive *bool      `json:"is_active"`
}

type UpdateUserRequest struct {
	Name     string     `json:"name" validate:"omitempty,min=2,max=100"`
	RoleID   *uuid.UUID `json:"role_id"`
	IsActive *bool      `json:"is_active"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
}

type UserUsecase interface {
	FindAll(filter repository.UserFilter) ([]entity.User, int64, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	Create(req CreateUserRequest) (*entity.User, error)
	Update(id uuid.UUID, req UpdateUserRequest) (*entity.User, error)
	Delete(id uuid.UUID) error
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*entity.User, error)
	UpdatePhoto(userID uuid.UUID, file *multipart.FileHeader) (*entity.User, error)
}
