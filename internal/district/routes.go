// package districts handles district-related routes and operations
package district

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/districts", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Get all districts
		r.Get("/", GetDistrictsHandler)

		// Get districts by state id
		r.Get(
			"/states/{state_id}",
			GetDistrictsByStateIdHandler,
		)

	})
}