package academy

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	// ==================================================
	// Academy Routes
	// ==================================================

	r.Route("/academies", func(r chi.Router) {

		// --------------------------------------------------
		// Auth Required For All Academy Routes
		// --------------------------------------------------

		r.Use(middleware.AuthMiddleware)

		// ==================================================
		// District Admin Routes
		// ==================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"district_admin",
			))

			r.Post(
				"/",
				CreateAcademyHandler,
			)

			r.Get(
				"/my-district",
				GetDistrictAdminAcademiesHandler,
			)
		})

		// ==================================================
		// Super Admin Routes
		// ==================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"super_admin",
			))

			r.Get(
				"/",
				GetAcademiesHandler,
			)
		})

		// ==================================================
		// Player Management Routes
		// ==================================================

		r.Route("/players", func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))

			r.Get(
				"/",
				GetAcademyPlayersHandler,
			)

			r.Get(
				"/{playerID}",
				GetAcademyPlayerHandler,
			)

			r.Post(
				"/{playerID}/assign-coach",
				AssignCoachHandler,
			)

			r.Delete(
				"/{playerID}/coach",
				RemoveCoachHandler,
			)
		})

		// ==================================================
		// Coach Management Routes
		// ==================================================

		r.Route("/coaches", func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))

			r.Get(
				"/",
				GetAcademyCoachesHandler,
			)

			r.Get(
				"/{coachID}",
				GetAcademyCoachHandler,
			)
		})

		// ==================================================
		// Building Management Routes
		// ==================================================

		r.Route("/buildings", func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
			))

			// ==============================================
			// Building CRUD
			// ==============================================

			r.Post(
				"/",
				CreateAcademyBuildingHandler,
			)

			r.Get(
				"/",
				GetAcademyBuildingsHandler,
			)

			r.Get(
				"/{buildingID}",
				GetAcademyBuildingHandler,
			)

			r.Patch(
				"/{buildingID}",
				UpdateAcademyBuildingHandler,
			)

			// ==============================================
			// Building Disciplines
			// ==============================================

			r.Post(
				"/{buildingID}/disciplines",
				AddAcademyBuildingDisciplineHandler,
			)

			r.Delete(
				"/{buildingID}/disciplines/{disciplineID}",
				RemoveAcademyBuildingDisciplineHandler,
			)

			// ==============================================
			// Building Events
			// ==============================================

			r.Post(
				"/{buildingID}/events",
				AddAcademyBuildingEventHandler,
			)

			r.Delete(
				"/{buildingID}/events/{eventID}",
				RemoveAcademyBuildingEventHandler,
			)

			// ==============================================
			// Building Lanes
			// ==============================================

			r.Post(
				"/{buildingID}/lanes",
				AddAcademyBuildingLaneHandler,
			)

			r.Get(
				"/{buildingID}/lanes",
				GetAcademyBuildingLanesHandler,
			)

			r.Patch(
				"/lanes/{laneID}",
				UpdateAcademyBuildingLaneHandler,
			)

			r.Delete(
				"/lanes/{laneID}",
				DeleteAcademyBuildingLaneHandler,
			)
		})

		// ==================================================
		// Lane Availability
		// ==================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.RequireRole(
				"academy_admin",
				"player",
			))

			r.Get(
				"/buildings/{buildingID}/lanes/available",
				GetAvailableLanesHandler,
			)
		})
	})
}