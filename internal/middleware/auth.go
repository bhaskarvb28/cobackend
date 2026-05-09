package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cobackend/internal/auth"
	"cobackend/internal/shared"
	"cobackend/internal/utils"

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

	if len(jwtSecret) == 0 {
		panic("JWT_SECRET is not configured")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: "Missing authorization header",
				},
			)

			return
		}

		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 || splitToken[0] != "Bearer" {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: "Invalid authorization header",
				},
			)

			return
		}

		tokenString := splitToken[1]

		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: "Invalid token",
				},
			)

			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims.Subject,
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