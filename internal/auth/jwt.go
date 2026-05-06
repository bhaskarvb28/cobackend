package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"sub"`
	RoleID string `json:"role"`
	Role string `json:"role_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string, roleID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}