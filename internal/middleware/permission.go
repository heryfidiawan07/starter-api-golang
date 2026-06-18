package middleware

import (
	"starter-api-golang/internal/domain/repository"
	"starter-api-golang/pkg/response"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permissionName string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsRoot(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		user, err := userRepo.FindByID(userID)
		if err != nil {
			response.Forbidden(c, "access denied")
			c.Abort()
			return
		}

		if user.Role == nil {
			response.Forbidden(c, "access denied: no role assigned")
			c.Abort()
			return
		}

		for _, perm := range user.Role.Permissions {
			if perm.Name == permissionName {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "access denied: insufficient permissions")
		c.Abort()
	}
}
