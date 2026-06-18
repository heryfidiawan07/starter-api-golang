package usecase

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type PermissionUsecase interface {
	FindAll() ([]entity.Permission, error)
	FindTree() ([]entity.Permission, error)
	FindByRoleID(roleID uuid.UUID) ([]entity.Permission, error)
}
