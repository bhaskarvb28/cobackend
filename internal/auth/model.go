package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Internal auth/database model
type AuthUser struct {
	ID       string
	Email    string
	Password string
	RoleID   string
	Role     string
}

// JWT claims
type Claims struct {
	RoleID   string `json:"role"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

// Public user response
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Login API response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

