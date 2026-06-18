package repository

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type SocialAccountRepository interface {
	FindByProviderAndID(provider entity.SocialProvider, providerID string) (*entity.SocialAccount, error)
	FindByUserID(userID uuid.UUID) ([]entity.SocialAccount, error)
	Create(account *entity.SocialAccount) error
	Update(account *entity.SocialAccount) error
}
