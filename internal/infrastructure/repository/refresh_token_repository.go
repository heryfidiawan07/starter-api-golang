package repository

import (
	"time"

	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) domainRepo.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) FindByToken(token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.db.First(&rt, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepository) FindActiveByUserID(userID uuid.UUID) ([]entity.RefreshToken, error) {
	var tokens []entity.RefreshToken
	err := r.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Find(&tokens).Error
	return tokens, err
}

func (r *refreshTokenRepository) Create(token *entity.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) Revoke(token string) error {
	now := time.Now()
	return r.db.Model(&entity.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", now).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&entity.RefreshToken{}).Error
}
