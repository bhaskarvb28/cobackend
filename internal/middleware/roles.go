package middleware

import "net/http"

func RequireRole(roles ...string) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleNameKey).(string)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range roles {

				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)

		})
	}
}