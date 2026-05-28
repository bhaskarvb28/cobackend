package session

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/practice-session", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		// ------------------------------------------------------------------
		// Player Admin Routes
		// ------------------------------------------------------------------

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"player",
			))

			// r.Get(
			// 	"/shooting-events",
			// 	GetAvailableShootingEventsHandler,
			// )

			// r.Get(
			// 	"/buildings",
			// 	GetCompatibleBuildingsHandler,
			// )

			// routes.go
			r.Post(
				"/",
				StartPracticeSessionHandler,
			)			
		})
	})
}