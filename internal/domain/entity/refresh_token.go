package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Token     string     `gorm:"not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `gorm:"default:null" json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *RefreshToken) IsValid() bool {
	return r.RevokedAt == nil && r.ExpiresAt.After(time.Now())
}
