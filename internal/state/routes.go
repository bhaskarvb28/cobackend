// Package state handles state-related routes and operations.
package state

import (
	"github.com/go-chi/chi/v5"

	"cobackend/internal/district"
	"cobackend/internal/middleware"
)

// RegisterRoutes registers all state-related routes.
//
// Routes:
//
//	Protected:
//		GET    /states
//		GET    /states/{state_id}/districts
//
// Middleware:
//
//	- AuthMiddleware
func RegisterRoutes(r chi.Router) {
	r.Route("/states", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		// Get all states.
		r.Get("/", GetStatesHandler)

		// Get all districts belonging to a state.
		r.Get(
			"/{state_id}/districts",
			district.GetDistrictsByStateIdHandler,
		)
	})
}