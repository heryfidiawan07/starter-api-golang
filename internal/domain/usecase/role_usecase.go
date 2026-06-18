package usecase

import (
	"starter-api-golang/internal/domain/entity"
	"starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name          string      `json:"name" validate:"required,min=2,max=100"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string      `json:"name" validate:"omitempty,min=2,max=100"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

type RoleUsecase interface {
	FindAll(filter repository.RoleFilter) ([]entity.Role, int64, error)
	FindByID(id uuid.UUID) (*entity.Role, error)
	Create(req CreateRoleRequest) (*entity.Role, error)
	Update(id uuid.UUID, req UpdateRoleRequest) (*entity.Role, error)
	Delete(id uuid.UUID) error
}
