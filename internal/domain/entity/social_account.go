package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialProvider string

const (
	ProviderGoogle   SocialProvider = "google"
	ProviderFacebook SocialProvider = "facebook"
)

type SocialAccount struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Provider      SocialProvider `gorm:"not null" json:"provider"`
	ProviderID    string         `gorm:"not null" json:"provider_id"`
	ProviderEmail string         `gorm:"not null" json:"provider_email"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *SocialAccount) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
