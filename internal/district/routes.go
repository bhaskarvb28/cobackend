// package districts handles district-related routes and operations 
package district

import (
	"github.com/go-chi/chi/v5"
	"cobackend/internal/middleware"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/districts", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Get all districts
		r.Get("/", GetDistrictsHandler)
	})
}