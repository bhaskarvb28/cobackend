package academyAdmin

import (
	"github.com/go-chi/chi/v5"
	"cobackend/internal/middleware"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/academy-admins", func (r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Group(func (r chi.Router) {

			r.Use(middleware.RequireRole(
				"district_admin",
			))
		})

		r.Post("/invite", InviteAcademyAdminHandler)
		r.Get("/", GetAcademyAdminsHandler)
		r.Get("/{profile_id}", GetAcademyAdminByIDHandler)

	})
}
