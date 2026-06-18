package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionType string

const (
	PermissionTypeCategory PermissionType = "category"
	PermissionTypeMenu     PermissionType = "menu"
	PermissionTypeAction   PermissionType = "action"
)

type Permission struct {
	ID       uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	ParentID *uuid.UUID     `gorm:"type:char(36);default:null" json:"parent_id"`
	Parent   *Permission    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Permission   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Name     string         `gorm:"uniqueIndex;not null" json:"name"`
	Label    string         `gorm:"not null" json:"label"`
	Type     PermissionType `gorm:"not null" json:"type"`
	Route    *string        `gorm:"default:null" json:"route"`
	Icon     *string        `gorm:"default:null" json:"icon"`
	Order    int            `gorm:"default:0" json:"order"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
