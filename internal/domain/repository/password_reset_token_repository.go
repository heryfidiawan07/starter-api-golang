package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type PasswordResetTokenRepository interface {
	FindByToken(token string) (*entity.PasswordResetToken, error)
	FindActiveByUserID(userID uuid.UUID) (*entity.PasswordResetToken, error)
	Create(token *entity.PasswordResetToken) error
	MarkUsed(token string) error
	DeleteByUserID(userID uuid.UUID) error
}
