package repository

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domainRepo.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll(filter domainRepo.RoleFilter) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	query := r.db.Model(&entity.Role{}).Preload("Permissions")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	query.Count(&total)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) FindByID(id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *entity.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Role{}, "id = ?", id).Error
}

func (r *roleRepository) SyncPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&entity.RolePermission{}, "role_id = ?", roleID).Error; err != nil {
			return err
		}
		for _, permID := range permissionIDs {
			rp := entity.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
