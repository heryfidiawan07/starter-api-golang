package main

import (
	"log"

	"starter-api-golang/internal/config"
	"starter-api-golang/internal/database"
	"starter-api-golang/internal/handler"
	infraRepo "starter-api-golang/internal/infrastructure/repository"
	"starter-api-golang/internal/router"
	"starter-api-golang/internal/service"
	"starter-api-golang/pkg/email"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("Database seeding failed: %v", err)
	}

	// Repositories
	userRepo := infraRepo.NewUserRepository(db)
	roleRepo := infraRepo.NewRoleRepository(db)
	permissionRepo := infraRepo.NewPermissionRepository(db)
	socialRepo := infraRepo.NewSocialAccountRepository(db)
	refreshTokenRepo := infraRepo.NewRefreshTokenRepository(db)
	resetTokenRepo := infraRepo.NewPasswordResetTokenRepository(db)

	// Mailer
	mailer := email.New(
		cfg.Mail.Host,
		cfg.Mail.Port,
		cfg.Mail.User,
		cfg.Mail.Pass,
		cfg.Mail.From,
		cfg.Mail.FromName,
	)

	// Services
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, resetTokenRepo, socialRepo, mailer, cfg)
	userSvc := service.NewUserService(userRepo, cfg)
	roleSvc := service.NewRoleService(roleRepo, permissionRepo, userRepo)
	permissionSvc := service.NewPermissionService(permissionRepo)

	// Handlers
	authHandler       := handler.NewAuthHandler(authSvc)
	userHandler       := handler.NewUserHandler(userSvc)
	roleHandler       := handler.NewRoleHandler(roleSvc)
	permissionHandler := handler.NewPermissionHandler(permissionSvc)
	lookupHandler     := handler.NewLookupHandler(roleRepo, permissionRepo)

	// Router
	r := router.New(cfg, userRepo, authHandler, userHandler, roleHandler, permissionHandler, lookupHandler)
	engine := r.Setup()

	addr := ":" + cfg.App.Port
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
