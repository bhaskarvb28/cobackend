package role

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/roles", func(r chi.Router) {

		r.Use(
			middleware.AuthMiddleware,
		)

		r.Get(
			"/invitable",
			GetInvitableRoleOptionsHandler,
		)
	})
}