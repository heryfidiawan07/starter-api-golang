package repository

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domainRepo.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll() ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Order("`order` ASC").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByID(id uuid.UUID) (*entity.Permission, error) {
	var p entity.Permission
	err := r.db.First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *permissionRepository) FindByName(name string) (*entity.Permission, error) {
	var p entity.Permission
	err := r.db.First(&p, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *permissionRepository) FindTree() ([]entity.Permission, error) {
	var all []entity.Permission
	if err := r.db.Order("`order` ASC").Find(&all).Error; err != nil {
		return nil, err
	}

	// Group direct children by parent ID
	childrenOf := make(map[uuid.UUID][]entity.Permission)
	for _, p := range all {
		if p.ParentID != nil {
			childrenOf[*p.ParentID] = append(childrenOf[*p.ParentID], p)
		}
	}

	// Recursively attach children before returning each node
	var attachChildren func(p entity.Permission) entity.Permission
	attachChildren = func(p entity.Permission) entity.Permission {
		if kids, ok := childrenOf[p.ID]; ok {
			for _, kid := range kids {
				p.Children = append(p.Children, attachChildren(kid))
			}
		}
		return p
	}

	var roots []entity.Permission
	for _, p := range all {
		if p.ParentID == nil {
			roots = append(roots, attachChildren(p))
		}
	}
	return roots, nil
}

func (r *permissionRepository) FindByRoleID(roleID uuid.UUID) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) Create(permission *entity.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *entity.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Permission{}, "id = ?", id).Error
}
