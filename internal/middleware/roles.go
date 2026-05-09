package middleware

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(RoleNameKey).(string)

			if !ok {
				utils.WriteJSON(w, http.StatusUnauthorized, shared.APIResponse{
					Success: false,
					Message: "unauthorized",
				})
				return
			}

			for _, allowedRole := range roles {

				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.WriteJSON(w, http.StatusForbidden, shared.APIResponse{
				Success: false,
				Message: "forbidden",
			})
		})
	}
}