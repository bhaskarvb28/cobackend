package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, roleID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		RoleName: role,
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