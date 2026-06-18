package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Token     string     `gorm:"not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `gorm:"default:null" json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *PasswordResetToken) IsValid() bool {
	return p.UsedAt == nil && p.ExpiresAt.After(time.Now())
}
