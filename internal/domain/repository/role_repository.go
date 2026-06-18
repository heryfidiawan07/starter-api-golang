package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type RoleFilter struct {
	Search  string
	Page    int
	PerPage int
}

type RoleRepository interface {
	FindAll(filter RoleFilter) ([]entity.Role, int64, error)
	FindByID(id uuid.UUID) (*entity.Role, error)
	FindByName(name string) (*entity.Role, error)
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uuid.UUID) error
	SyncPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
}
