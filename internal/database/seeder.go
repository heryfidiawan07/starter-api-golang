package database

import (
	"log"
	"time"

	"starter-api-golang/internal/domain/entity"
	"starter-api-golang/pkg/hash"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Running database seeders...")

	if err := seedPermissions(db); err != nil {
		return err
	}

	if err := seedRootUser(db); err != nil {
		return err
	}

	log.Println("Seeders completed successfully")
	return nil
}

func seedPermissions(db *gorm.DB) error {
	var count int64
	db.Model(&entity.Permission{}).Count(&count)
	if count > 0 {
		return nil
	}

	mainCat := entity.Permission{
		ID:    uuid.New(),
		Name:  "main",
		Label: "Main",
		Type:  entity.PermissionTypeCategory,
		Order: 1,
	}

	settingsCat := entity.Permission{
		ID:    uuid.New(),
		Name:  "settings",
		Label: "Settings",
		Type:  entity.PermissionTypeCategory,
		Order: 2,
	}

	if err := db.Create(&mainCat).Error; err != nil {
		return err
	}
	if err := db.Create(&settingsCat).Error; err != nil {
		return err
	}

	dashRoute := "/dashboard"
	dashIcon := "layout-dashboard"
	dashboard := entity.Permission{
		ID:       uuid.New(),
		ParentID: &mainCat.ID,
		Name:     "dashboard:index",
		Label:    "Dashboard",
		Type:     entity.PermissionTypeMenu,
		Route:    &dashRoute,
		Icon:     &dashIcon,
		Order:    1,
	}
	if err := db.Create(&dashboard).Error; err != nil {
		return err
	}

	adminRoute := ""
	adminIcon := "shield"
	adminMenu := entity.Permission{
		ID:       uuid.New(),
		ParentID: &settingsCat.ID,
		Name:     "administrator",
		Label:    "Administrator",
		Type:     entity.PermissionTypeMenu,
		Route:    &adminRoute,
		Icon:     &adminIcon,
		Order:    1,
	}
	if err := db.Create(&adminMenu).Error; err != nil {
		return err
	}

	userRoute := "/admin/users"
	userIcon := "users"
	userMenu := entity.Permission{
		ID:       uuid.New(),
		ParentID: &adminMenu.ID,
		Name:     "user:index",
		Label:    "User",
		Type:     entity.PermissionTypeMenu,
		Route:    &userRoute,
		Icon:     &userIcon,
		Order:    1,
	}
	if err := db.Create(&userMenu).Error; err != nil {
		return err
	}

	userActions := []entity.Permission{
		{ID: uuid.New(), ParentID: &userMenu.ID, Name: "user:create", Label: "Create User", Type: entity.PermissionTypeAction, Order: 1},
		{ID: uuid.New(), ParentID: &userMenu.ID, Name: "user:edit", Label: "Edit User", Type: entity.PermissionTypeAction, Order: 2},
		{ID: uuid.New(), ParentID: &userMenu.ID, Name: "user:delete", Label: "Delete User", Type: entity.PermissionTypeAction, Order: 3},
	}
	for _, a := range userActions {
		if err := db.Create(&a).Error; err != nil {
			return err
		}
	}

	roleRoute := "/admin/roles"
	roleIcon := "key"
	roleMenu := entity.Permission{
		ID:       uuid.New(),
		ParentID: &adminMenu.ID,
		Name:     "role:index",
		Label:    "Role",
		Type:     entity.PermissionTypeMenu,
		Route:    &roleRoute,
		Icon:     &roleIcon,
		Order:    2,
	}
	if err := db.Create(&roleMenu).Error; err != nil {
		return err
	}

	roleActions := []entity.Permission{
		{ID: uuid.New(), ParentID: &roleMenu.ID, Name: "role:create", Label: "Create Role", Type: entity.PermissionTypeAction, Order: 1},
		{ID: uuid.New(), ParentID: &roleMenu.ID, Name: "role:edit", Label: "Edit Role", Type: entity.PermissionTypeAction, Order: 2},
		{ID: uuid.New(), ParentID: &roleMenu.ID, Name: "role:delete", Label: "Delete Role", Type: entity.PermissionTypeAction, Order: 3},
	}
	for _, a := range roleActions {
		if err := db.Create(&a).Error; err != nil {
			return err
		}
	}

	permRoute := "/admin/permissions"
	permIcon := "lock"
	permMenu := entity.Permission{
		ID:       uuid.New(),
		ParentID: &adminMenu.ID,
		Name:     "permission:index",
		Label:    "Permission",
		Type:     entity.PermissionTypeMenu,
		Route:    &permRoute,
		Icon:     &permIcon,
		Order:    3,
	}
	if err := db.Create(&permMenu).Error; err != nil {
		return err
	}

	return nil
}

func seedRootUser(db *gorm.DB) error {
	var count int64
	db.Model(&entity.User{}).Where("is_root = ?", true).Count(&count)
	if count > 0 {
		return nil
	}

	hashed, err := hash.Make("password")
	if err != nil {
		return err
	}

	now := time.Now()
	root := entity.User{
		Username:   "root",
		Name:       "Root Administrator",
		Email:      "root@example.com",
		Password:   &hashed,
		IsRoot:     true,
		Status:     true,
		VerifiedAt: &now,
	}

	if err := db.Create(&root).Error; err != nil {
		return err
	}

	log.Printf("Root user created: email=root@example.com password=password")
	return nil
}
