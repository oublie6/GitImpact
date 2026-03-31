package middleware

import (
	"strings"

	"gitimpact/backend/pkg/auth"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			response.Err(c, 401, "missing token")
			c.Abort()
			return
		}
		claims, err := auth.ParseToken(secret, strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			response.Err(c, 401, "invalid token")
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
