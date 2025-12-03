package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

var (
	ErrInvalidToken = errors.New("Invalid token")
	ErrInvalidClaim = errors.New("Invalid token claim")
)

var secretKey = []byte("qs-lzh-security-key")

func CreateToken(username string, userRole model.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
			"username":  username,
			"user_role": userRole,
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaim
	}
	return claims, nil
}
