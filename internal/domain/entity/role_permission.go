package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	RoleID       uuid.UUID  `gorm:"type:char(36);not null;index" json:"role_id"`
	PermissionID uuid.UUID  `gorm:"type:char(36);not null;index" json:"permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"-"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (rp *RolePermission) BeforeCreate(tx *gorm.DB) error {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	return nil
}
