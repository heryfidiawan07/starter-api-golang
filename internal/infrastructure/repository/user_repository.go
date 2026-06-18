package repository

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(filter domainRepo.UserFilter) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.Model(&entity.User{}).Preload("Role")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ?", like, like, like)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.RoleID != nil {
		query = query.Where("role_id = ?", *filter.RoleID)
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
	err := query.Offset(offset).Limit(perPage).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role.Permissions").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role.Permissions").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}
