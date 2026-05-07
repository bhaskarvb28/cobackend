package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cobackend/internal/auth"

	"github.com/golang-jwt/jwt/v5"

)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	RoleIDKey   contextKey = "role_id"
	RoleNameKey contextKey = "role_name"
)

func AuthMiddleware(next http.Handler) http.Handler {

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims.UserID,
		)

		ctx = context.WithValue(
			ctx,
			RoleIDKey,
			claims.RoleID,
		)

		ctx = context.WithValue(
			ctx,
			RoleNameKey,
			claims.RoleName,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}