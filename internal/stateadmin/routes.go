package stateadmin

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/state-admins", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Group(func(r chi.Router) {

			r.Use(
				middleware.RequireRole("super_admin"),
			)

			r.Post("/", CreateStateAdminHandler)

			// Get State Admin Routes
			r.Get("/", GetStateAdminsHandler)
			
			r.Put("/{id}/assigned-state", UpdateAssignedStateHandler)
			r.Delete("/{id}", DeleteStateAdminHandler)
		})
	})
}