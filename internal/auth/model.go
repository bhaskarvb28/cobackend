package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// type RegisterInput struct {
// 	FirstName     string `json:"first_name"`
// 	LastName      string `json:"last_name"`
// 	Email         string `json:"email"`
// 	Password      string `json:"password"` // plain password
// 	ContactNumber string `json:"contact_number"`
// }

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID       string
	Email    string
	Password string
	RoleID   string
	Role     string
}

type Claims struct {
	UserID string `json:"sub"`
	RoleID string `json:"role"`
	Role   string `json:"role_name"`
	jwt.RegisteredClaims
}
