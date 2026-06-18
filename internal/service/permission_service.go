package service

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"

	"github.com/google/uuid"
)

type permissionService struct {
	permissionRepo domainRepo.PermissionRepository
}

func NewPermissionService(permissionRepo domainRepo.PermissionRepository) usecase.PermissionUsecase {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) FindAll() ([]entity.Permission, error) {
	return s.permissionRepo.FindAll()
}

func (s *permissionService) FindTree() ([]entity.Permission, error) {
	return s.permissionRepo.FindTree()
}

func (s *permissionService) FindByRoleID(roleID uuid.UUID) ([]entity.Permission, error) {
	return s.permissionRepo.FindByRoleID(roleID)
}
