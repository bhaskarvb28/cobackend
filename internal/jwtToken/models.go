package jwtToken

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWT claims
type Claims struct {
	RoleID   int `json:"role"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}