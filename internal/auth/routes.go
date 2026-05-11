package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", LoginHandler)
		r.Post("/accept-invitation", AcceptInviteHandler)
	})
}