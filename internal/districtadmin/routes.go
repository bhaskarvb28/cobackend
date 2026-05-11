package districtadmin

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/district-admins", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(
				"state_admin",
			))

			r.Post("/invite", InviteDistrictAdminHandler)

			// r.Post("/", CreateDistrictAdminHandler)
			// r.Get("/", GetDistrictAdminsHandler)
			// r.Put("/{id}", UpdateDistrictAdminHandler)
			// r.Delete("/{id}", DeleteDistrictAdminHandler)
		})
	})
}
