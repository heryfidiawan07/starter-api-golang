package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Username   string         `gorm:"uniqueIndex;not null" json:"username"`
	Name       string         `gorm:"not null" json:"name"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Password   *string        `gorm:"default:null" json:"-"`
	Photo      *string        `gorm:"default:null" json:"photo"`
	VerifiedAt *time.Time     `gorm:"default:null" json:"verified_at"`
	IsRoot     bool           `gorm:"default:false" json:"is_root"`
	RoleID     *uuid.UUID     `gorm:"type:char(36);default:null" json:"role_id"`
	Role       *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Status     bool           `gorm:"default:true" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
