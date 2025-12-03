package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qs-lzh/movie-reservation/internal/dto"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		if !exists {
			dto.Forbidden(c, "Not permitted to use")
			return
		}
		claims, ok := claimsAny.(jwt.MapClaims)
		if !ok {
			dto.InternalServerError(c, "Failed to parse token claims")
			return
		}

		userRole, ok := claims["user_role"]
		if !ok {
			dto.InternalServerError(c, "Failed to get user role from claims")
			return
		}
		if userRole != "admin" {
			dto.Forbidden(c, "Not permitted to use")
			return
		}

		c.Next()

	}
}
