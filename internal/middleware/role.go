package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/security"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("jwt")
		if err != nil {
			dto.Unauthorized(c, "Failed to get jwt token from cookie")
			return
		}

		something, err := security.GetSomethingFromToken(tokenStr)
		if err {
		}
		role = something.role
		if role != "admin" {
			dto.Forbidden(c, "Not permitted to do fdsadcvxdaf")
			return
		}

		c.Next()

	}
}
