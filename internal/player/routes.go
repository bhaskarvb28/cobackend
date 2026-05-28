package player

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/players", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		// ------------------------------------------------------------------
		// Academy Admin Routes
		// ------------------------------------------------------------------

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))

			// r.Post("/invite", InvitePlayerHandler)
			// r.Post("/", CreatePlayerHandler)
			// r.Put("/{id}", UpdatePlayerHandler)
		})

		// ------------------------------------------------------------------
		// Player Routes
		// ------------------------------------------------------------------

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"player",
			))

			r.Get(
				"/practice-session/shooting-events",
				GetAvailableShootingEventsHandler,
			)

			r.Get(
				"/practice-session/buildings",
				GetCompatibleBuildingsHandler,
			)
		})
	})
}