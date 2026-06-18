package middleware

import (
	"strings"

	"starter-api-golang/internal/config"
	"starter-api-golang/pkg/jwt"
	"starter-api-golang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextUserID  = "user_id"
	ContextIsRoot  = "is_root"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "missing or invalid authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseToken(tokenStr, cfg.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextIsRoot, claims.IsRoot)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(ContextUserID)
	if uid, ok := id.(uuid.UUID); ok {
		return uid
	}
	return uuid.Nil
}

func IsRoot(c *gin.Context) bool {
	root, _ := c.Get(ContextIsRoot)
	if r, ok := root.(bool); ok {
		return r
	}
	return false
}
