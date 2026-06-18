package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type PermissionRepository interface {
	FindAll() ([]entity.Permission, error)
	FindByID(id uuid.UUID) (*entity.Permission, error)
	FindByName(name string) (*entity.Permission, error)
	FindTree() ([]entity.Permission, error)
	FindByRoleID(roleID uuid.UUID) ([]entity.Permission, error)
	Create(permission *entity.Permission) error
	Update(permission *entity.Permission) error
	Delete(id uuid.UUID) error
}
