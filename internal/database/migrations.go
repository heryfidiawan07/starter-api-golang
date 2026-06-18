package database

import (
	"log"

	"starter-api-golang/internal/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&entity.Permission{},
		&entity.Role{},
		&entity.User{},
		&entity.RolePermission{},
		&entity.SocialAccount{},
		&entity.RefreshToken{},
		&entity.PasswordResetToken{},
	)
	if err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
