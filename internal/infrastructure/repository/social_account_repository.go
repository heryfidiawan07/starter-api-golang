package repository

import (
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type socialAccountRepository struct {
	db *gorm.DB
}

func NewSocialAccountRepository(db *gorm.DB) domainRepo.SocialAccountRepository {
	return &socialAccountRepository{db: db}
}

func (r *socialAccountRepository) FindByProviderAndID(provider entity.SocialProvider, providerID string) (*entity.SocialAccount, error) {
	var account entity.SocialAccount
	err := r.db.First(&account, "provider = ? AND provider_id = ?", provider, providerID).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *socialAccountRepository) FindByUserID(userID uuid.UUID) ([]entity.SocialAccount, error) {
	var accounts []entity.SocialAccount
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *socialAccountRepository) Create(account *entity.SocialAccount) error {
	return r.db.Create(account).Error
}

func (r *socialAccountRepository) Update(account *entity.SocialAccount) error {
	return r.db.Save(account).Error
}
