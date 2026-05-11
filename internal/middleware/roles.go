package middleware

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().
				Value(RoleNameKey).
				( string)

			if !ok {

				utils.WriteError(
					w,
					shared.NewAPIError(
						http.StatusUnauthorized,
						"unauthorized",
					),
					"authorization failed",
				)

				return
			}

			for _, allowedRole := range roles {

				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.WriteError(
				w,
				shared.NewAPIError(
					http.StatusForbidden,
					"forbidden",
				),
				"authorization failed",
			)
		})
	}
}