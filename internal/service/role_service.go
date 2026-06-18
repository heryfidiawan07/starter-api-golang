package service

import (
	"errors"

	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound = errors.New("role not found")
	ErrRoleNameTaken = errors.New("role name already taken")
	ErrRoleInUse     = errors.New("role is assigned to users and cannot be deleted")
)

type roleService struct {
	roleRepo       domainRepo.RoleRepository
	permissionRepo domainRepo.PermissionRepository
	userRepo       domainRepo.UserRepository
}

func NewRoleService(
	roleRepo domainRepo.RoleRepository,
	permissionRepo domainRepo.PermissionRepository,
	userRepo domainRepo.UserRepository,
) usecase.RoleUsecase {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *roleService) FindAll(filter domainRepo.RoleFilter) ([]entity.Role, int64, error) {
	return s.roleRepo.FindAll(filter)
}

func (s *roleService) FindByID(id uuid.UUID) (*entity.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	}
	return role, err
}

func (s *roleService) Create(req usecase.CreateRoleRequest) (*entity.Role, error) {
	if _, err := s.roleRepo.FindByName(req.Name); err == nil {
		return nil, ErrRoleNameTaken
	}

	desc := req.Description
	role := &entity.Role{
		Name:        req.Name,
		Description: &desc,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepo.SyncPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return s.roleRepo.FindByID(role.ID)
}

func (s *roleService) Update(id uuid.UUID, req usecase.UpdateRoleRequest) (*entity.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}

	if req.Name != "" && req.Name != role.Name {
		if _, err := s.roleRepo.FindByName(req.Name); err == nil {
			return nil, ErrRoleNameTaken
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = &req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	if err := s.roleRepo.SyncPermissions(role.ID, req.PermissionIDs); err != nil {
		return nil, err
	}

	return s.roleRepo.FindByID(role.ID)
}

func (s *roleService) Delete(id uuid.UUID) error {
	if _, err := s.roleRepo.FindByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRoleNotFound
	}

	users, _, err := s.userRepo.FindAll(domainRepo.UserFilter{RoleID: &id, Page: 1, PerPage: 1})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return ErrRoleInUse
	}

	return s.roleRepo.Delete(id)
}
