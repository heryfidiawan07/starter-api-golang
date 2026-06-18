package repository

import (
	"time"

	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) domainRepo.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) FindByToken(token string) (*entity.PasswordResetToken, error) {
	var prt entity.PasswordResetToken
	err := r.db.First(&prt, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &prt, nil
}

func (r *passwordResetTokenRepository) FindActiveByUserID(userID uuid.UUID) (*entity.PasswordResetToken, error) {
	var prt entity.PasswordResetToken
	err := r.db.Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, time.Now()).
		First(&prt).Error
	if err != nil {
		return nil, err
	}
	return &prt, nil
}

func (r *passwordResetTokenRepository) Create(token *entity.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetTokenRepository) MarkUsed(token string) error {
	now := time.Now()
	return r.db.Model(&entity.PasswordResetToken{}).
		Where("token = ?", token).
		Update("used_at", now).Error
}

func (r *passwordResetTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.PasswordResetToken{}).Error
}
