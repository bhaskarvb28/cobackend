package academy

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/academies", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"district_admin",
			))

			r.Post("/", CreateAcademyHandler)
			r.Get("/", GetAcademiesHandler)

		})

		r.Route("/buildings", func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))

			r.Post("/", CreateAcademyBuildingHandler)
			r.Get("/", GetAcademyBuildingsHandler)

			r.Post(
				"/{buildingID}/disciplines",
				AddAcademyBuildingDisciplineHandler,
			)

			r.Post(
				"/{buildingID}/events",
				AddAcademyBuildingEventHandler,
			)
		})
	})
}