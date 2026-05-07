package superadmin

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/super-admin", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Route("/state-admins", func(r chi.Router) {

			// SUPER ADMIN ONLY
			r.With(
				middleware.RequireRole("super_admin"),
			).Post("/", CreateStateAdminHandler)

			// SUPER ADMIN ONLY
			r.With(
				middleware.RequireRole("super_admin"),
			).Get("/", GetStateAdminsHandler)

			// SUPER ADMIN + STATE ADMIN
			r.With(
				middleware.RequireRole(
					"super_admin",
				),
			).Put("/{id}/assigned-state", UpdateAssignedStateHandler)

			// SUPER ADMIN ONLY
			r.With(
				middleware.RequireRole("super_admin"),
			).Delete("/{id}", DeleteStateAdminHandler)



			// SUPER ADMIN + STATE ADMIN
			// r.With(
			// 	middleware.RequireRole(
			// 		"super_admin",
			// 		"state_admin",
			// 	),
			// ).Get("/{id}", GetStateAdminHandler)


		})
	})
}