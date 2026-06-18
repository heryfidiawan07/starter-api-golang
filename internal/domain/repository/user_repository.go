package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type UserFilter struct {
	Search  string
	Status  *bool
	RoleID  *uuid.UUID
	Page    int
	PerPage int
}

type UserRepository interface {
	FindAll(filter UserFilter) ([]entity.User, int64, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uuid.UUID) error
}
