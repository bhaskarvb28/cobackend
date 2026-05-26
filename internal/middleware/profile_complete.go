package middleware

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

// RequireCompletedProfile blocks users
// whose onboarding/profile is incomplete.
//
// This middleware must be used after
// AuthMiddleware because it depends
// on profile_completed being injected
// into request context.
func RequireCompletedProfile(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		// ----------------------------------------------------------
		// Extract Profile Completion Status
		// ----------------------------------------------------------

		profileCompleted, ok := r.Context().
			Value(ProfileCompletedKey).
			(bool)

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
		// Block Incomplete Profiles
		// ----------------------------------------------------------

		if !profileCompleted {

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: "profile completion required",
				},
			)

			return
		}

		// ----------------------------------------------------------
		// Continue Request Chain
		// ----------------------------------------------------------

		next.ServeHTTP(
			w,
			r,
		)
	})
}