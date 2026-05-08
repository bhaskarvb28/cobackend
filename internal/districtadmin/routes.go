package districtadmin

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/state-admin", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Route("/district-admins", func(r chi.Router) {

			// STATE ADMIN + SUPER ADMIN: Create district admin
			r.With(
				middleware.RequireRole("super_admin", "state_admin"),
			).Post("/", CreateDistrictAdminHandler)

			// STATE ADMIN + SUPER ADMIN: Get all district admins
			r.With(
				middleware.RequireRole("super_admin", "state_admin"),
			).Get("/", GetDistrictAdminsHandler)

			// STATE ADMIN + SUPER ADMIN: Update district admin
			r.With(
				middleware.RequireRole("super_admin", "state_admin"),
			).Put("/{id}", UpdateDistrictAdminHandler)

			// STATE ADMIN + SUPER ADMIN: Delete district admin
			r.With(
				middleware.RequireRole("super_admin", "state_admin"),
			).Delete("/{id}", DeleteDistrictAdminHandler)

		})
	})
}
