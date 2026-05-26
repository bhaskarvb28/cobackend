package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cobackend/internal/jwtToken"
	"cobackend/internal/shared"
	"cobackend/internal/utils"
	"cobackend/internal/auth"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	
	// UserIdKey stores the authenticated user's ID.
	UserIDKey   contextKey = "user_id"

	// RoleIDKey stores the authenticated user's role ID.
	RoleIDKey   contextKey = "role_id"

	// RoleNameKey stores the authenticated user's role name/code.
	RoleNameKey contextKey = "role_name"

	ProfileCompletedKey contextKey = "profile_completed"
)

// AuthMiddleware validates JWT authentication tokens
// and injects authenticated user information into
// the request context
//
// Expected Authorization header:
//
// Authorization: Bearer <token>
// On successful validation, the middleware stores:
// - user_id
// - role_id
// - role_name
//
// in the request context for downstream handlers
func AuthMiddleware(next http.Handler) http.Handler {

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if len(jwtSecret) == 0 {
		panic("JWT_SECRET is not configured")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ----------------------------------------------------------
		// Read Authorization Header
		// ----------------------------------------------------------

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: shared.ErrMissingAuthorizationHeader.Error(),
				},
			)

			return
		}

		// ----------------------------------------------------------
		// Validate Bearer Token Format
		// ----------------------------------------------------------

		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 || splitToken[0] != "Bearer" {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: shared.ErrInvalidAuthorizationHeader.Error(),
				},
			)

			return
		}

		tokenString := splitToken[1]

		// ----------------------------------------------------------
		// Parse JWT Claims
		// ----------------------------------------------------------

		claims := &jwtToken.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				// Ensure token uses HMAC signing method.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return jwtSecret, nil
			},
		)

		// ----------------------------------------------------------
		// Validate Token
		// ----------------------------------------------------------

		if err != nil || !token.Valid {

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: shared.ErrInvalidToken.Error(),
				},
			)

			return
		}

		profileCompleted, err := auth.GetProfileCompletedStatus(
			r.Context(),
			claims.Subject,
			claims.RoleName,
		)

		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusInternalServerError,
				shared.APIResponse{
					Success: false,
					Message: "failed to fetch profile status",
				},
			)

			return
		}

		// ----------------------------------------------------------
		// Inject Claims Into Context
		// ----------------------------------------------------------

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

		ctx = context.WithValue(
			ctx,
			ProfileCompletedKey,
			profileCompleted,
		)

		// ----------------------------------------------------------
		// Continue Request Chain
		// ----------------------------------------------------------

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}