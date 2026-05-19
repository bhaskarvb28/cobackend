package academyCoach

import (
	"github.com/go-chi/chi/v5"
	"cobackend/internal/middleware"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/academy-coach", func (r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Group(func (r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))
		})

		r.Post("/invite", InviteAcademyCoachHandler)
	})
}
