package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	FindByToken(token string) (*entity.RefreshToken, error)
	FindActiveByUserID(userID uuid.UUID) ([]entity.RefreshToken, error)
	Create(token *entity.RefreshToken) error
	Revoke(token string) error
	RevokeAllByUserID(userID uuid.UUID) error
	DeleteExpired() error
}
