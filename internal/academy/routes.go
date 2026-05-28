package academy

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/academies", func(r chi.Router) {

		// --------------------------------------------------
		// Auth Required For All Academy Routes
		// --------------------------------------------------

		r.Use(middleware.AuthMiddleware)

		// --------------------------------------------------
		// District Admin Routes
		// --------------------------------------------------

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"district_admin",
			))

			r.Post("/", CreateAcademyHandler)
			r.Get("/", GetAcademiesHandler)
		})

		// --------------------------------------------------
		// Building Routes
		// --------------------------------------------------

		r.Route("/buildings", func(r chi.Router) {

			// ==============================================
			// Academy Admin Only
			// ==============================================

			r.Group(func(r chi.Router) {

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

				r.Post(
					"/{buildingID}/lanes",
					AddAcademyBuildingLaneHandler,
				)
			})

			// ==============================================
			// Academy Admin + Player
			// ==============================================

			r.Group(func(r chi.Router) {

				r.Use(middleware.RequireRole(
					"academy_admin",
					"player",
				))

				r.Get(
					"/{buildingID}/lanes/available",
					GetAvailableLanesHandler,
				)
			})
		})
	})
}