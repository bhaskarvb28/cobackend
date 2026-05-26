package middleware

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

// RequireRole restricts access to authenticated users
// whose role matches one of the allowed roles.
//
// This middleware must be used after AuthMiddleware,
// since it depends on role information being present
// in the request context.
//
// Example:
//
//	r.Use(
//		RequireRole(
//			"super_admin",
//			"state_admin",
//		),
//	)
func RequireRole(
	roles ...string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			// ----------------------------------------------------------
			// Extract Role From Context
			// ----------------------------------------------------------

			role, ok := r.Context().
				Value(RoleNameKey).
				(string)

			if !ok {

				utils.WriteJSON(
					w,
					http.StatusUnauthorized,
					shared.APIResponse{
						Success: false,
						Message: "unauthorized",
					},
				)

				return
			}

			// ----------------------------------------------------------
			// Check Allowed Roles
			// ----------------------------------------------------------

			for _, allowedRole := range roles {

				if role == allowedRole {

					// --------------------------------------------------
					// Continue Request Chain
					// --------------------------------------------------

					next.ServeHTTP(
						w,
						r,
					)

					return
				}
			}

			// ----------------------------------------------------------
			// Access Denied
			// ----------------------------------------------------------

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: "forbidden",
				},
			)
		})
	}
}