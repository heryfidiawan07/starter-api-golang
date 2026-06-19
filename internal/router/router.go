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
	lookup     *handler.LookupHandler
}

func New(
	cfg *config.Config,
	userRepo domainRepo.UserRepository,
	auth *handler.AuthHandler,
	user *handler.UserHandler,
	role *handler.RoleHandler,
	permission *handler.PermissionHandler,
	lookup *handler.LookupHandler,
) *Router {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(corsMiddleware())

	engine.Static("/storage/photos", cfg.Storage.Path)

	return &Router{
		engine:     engine,
		cfg:        cfg,
		userRepo:   userRepo,
		auth:       auth,
		user:       user,
		role:       role,
		permission: permission,
		lookup:     lookup,
	}
}

// corsMiddleware handles CORS headers and OPTIONS preflight requests.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Starter API — Golang"})
	})

	api := r.engine.Group("/api/v1")

	// Public config — returns OAuth Client IDs for the frontend SDK
	api.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"google_client_id": r.cfg.Google.ClientID,
				"facebook_app_id":  r.cfg.Facebook.ClientID,
			},
		})
	})

	// Public auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", r.auth.Register)
		authGroup.POST("/login", r.auth.Login)
		authGroup.POST("/refresh", r.auth.RefreshToken)
		authGroup.POST("/revoke", r.auth.RevokeToken)
		authGroup.POST("/forgot-password", r.auth.ForgotPassword)
		authGroup.POST("/reset-password", r.auth.ResetPassword)
		authGroup.GET("/verify-email", r.auth.VerifyEmail)
		authGroup.POST("/oauth/google", r.auth.GoogleAuth)
		authGroup.POST("/oauth/facebook", r.auth.FacebookAuth)
	}

	// Protected auth routes
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.Auth(r.cfg))
	{
		authProtected.GET("/me", r.auth.Me)
		authProtected.POST("/logout", r.auth.Logout)
		authProtected.POST("/change-password", r.auth.ChangePassword)
	}

	// Profile routes
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
		users.GET("/:id", middleware.RequirePermission("user:show", r.userRepo), r.user.Show)
		users.PUT("/:id", middleware.RequirePermission("user:edit", r.userRepo), r.user.Update)
		users.DELETE("/:id", middleware.RequirePermission("user:delete", r.userRepo), r.user.Destroy)
		users.POST("/:id/photo", middleware.RequirePermission("user:edit", r.userRepo), r.user.UpdateUserPhoto)
	}

	// Role management routes
	roles := api.Group("/roles")
	roles.Use(middleware.Auth(r.cfg))
	{
		roles.GET("", middleware.RequirePermission("role:index", r.userRepo), r.role.Index)
		roles.POST("", middleware.RequirePermission("role:create", r.userRepo), r.role.Store)
		roles.GET("/:id", middleware.RequirePermission("role:show", r.userRepo), r.role.Show)
		roles.PUT("/:id", middleware.RequirePermission("role:edit", r.userRepo), r.role.Update)
		roles.DELETE("/:id", middleware.RequirePermission("role:delete", r.userRepo), r.role.Destroy)
	}

	// Permission routes
	permissions := api.Group("/permissions")
	permissions.Use(middleware.Auth(r.cfg))
	{
		permissions.GET("", middleware.RequirePermission("permission:index", r.userRepo), r.permission.Index)
		permissions.GET("/tree", middleware.RequirePermission("permission:index", r.userRepo), r.permission.Tree)
		permissions.GET("/by-role/:role_id", middleware.RequirePermission("permission:index", r.userRepo), r.permission.ByRole)
	}

	// Lookup routes — auth only, no specific permission required.
	// Used by form dropdowns (e.g. role select in Add User, permission tree in Add Role).
	lookupGroup := api.Group("/lookup")
	lookupGroup.Use(middleware.Auth(r.cfg))
	{
		lookupGroup.GET("/roles", r.lookup.Roles)
		lookupGroup.GET("/permissions", r.lookup.Permissions)
	}

	return r.engine
}
