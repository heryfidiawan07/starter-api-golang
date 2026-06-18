package router

import (
	"net/http"

	"starter-api-golang/internal/config"
	"starter-api-golang/internal/handler"
	"starter-api-golang/internal/middleware"
	domainRepo "starter-api-golang/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine     *gin.Engine
	cfg        *config.Config
	userRepo   domainRepo.UserRepository
	auth       *handler.AuthHandler
	user       *handler.UserHandler
	role       *handler.RoleHandler
	permission *handler.PermissionHandler
}

func New(
	cfg *config.Config,
	userRepo domainRepo.UserRepository,
	auth *handler.AuthHandler,
	user *handler.UserHandler,
	role *handler.RoleHandler,
	permission *handler.PermissionHandler,
) *Router {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.Static("/storage/photos", cfg.Storage.Path)

	return &Router{
		engine:     engine,
		cfg:        cfg,
		userRepo:   userRepo,
		auth:       auth,
		user:       user,
		role:       role,
		permission: permission,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Starter API — Golang"})
	})

	api := r.engine.Group("/api/v1")

	// Public auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", r.auth.Register)
		authGroup.POST("/login", r.auth.Login)
		authGroup.POST("/refresh-token", r.auth.RefreshToken)
		authGroup.POST("/forgot-password", r.auth.ForgotPassword)
		authGroup.POST("/reset-password", r.auth.ResetPassword)
		authGroup.GET("/verify-email", r.auth.VerifyEmail)
		authGroup.POST("/google", r.auth.GoogleAuth)
		authGroup.POST("/facebook", r.auth.FacebookAuth)
	}

	// Protected auth routes
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.Auth(r.cfg))
	{
		authProtected.GET("/me", r.auth.Me)
		authProtected.POST("/logout", r.auth.Logout)
		authProtected.POST("/revoke-token", r.auth.RevokeToken)
		authProtected.PUT("/change-password", r.auth.ChangePassword)
	}

	// Profile routes (authenticated user)
	profile := api.Group("/profile")
	profile.Use(middleware.Auth(r.cfg))
	{
		profile.PUT("", r.user.UpdateProfile)
		profile.POST("/photo", r.user.UpdatePhoto)
	}

	// User management routes
	users := api.Group("/users")
	users.Use(middleware.Auth(r.cfg))
	{
		users.GET("", middleware.RequirePermission("user:index", r.userRepo), r.user.Index)
		users.POST("", middleware.RequirePermission("user:create", r.userRepo), r.user.Store)
		users.GET("/:id", middleware.RequirePermission("user:index", r.userRepo), r.user.Show)
		users.PUT("/:id", middleware.RequirePermission("user:edit", r.userRepo), r.user.Update)
		users.DELETE("/:id", middleware.RequirePermission("user:delete", r.userRepo), r.user.Destroy)
	}

	// Role management routes
	roles := api.Group("/roles")
	roles.Use(middleware.Auth(r.cfg))
	{
		roles.GET("", middleware.RequirePermission("role:index", r.userRepo), r.role.Index)
		roles.POST("", middleware.RequirePermission("role:create", r.userRepo), r.role.Store)
		roles.GET("/:id", middleware.RequirePermission("role:index", r.userRepo), r.role.Show)
		roles.PUT("/:id", middleware.RequirePermission("role:edit", r.userRepo), r.role.Update)
		roles.DELETE("/:id", middleware.RequirePermission("role:delete", r.userRepo), r.role.Destroy)
	}

	// Permission routes (read-only, any authenticated user)
	permissions := api.Group("/permissions")
	permissions.Use(middleware.Auth(r.cfg))
	{
		permissions.GET("", r.permission.Index)
		permissions.GET("/tree", r.permission.Tree)
		permissions.GET("/role/:role_id", r.permission.ByRole)
	}

	return r.engine
}
