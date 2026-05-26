// Package auth handles authentication,
// invitation acceptance, and JWT validation.
package auth

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers authentication routes.
//
// Routes:
//   - POST /login
//   - POST /accept-invitation
func RegisterRoutes(r chi.Router) {

	r.Route("/auth", func(r chi.Router) {

		r.Post("/login", LoginHandler)
		// r.Post("/accept-invitation", AcceptInviteHandler)
	})
}