package disciplines

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(r chi.Router){
	r.Route("/disciplines", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Use(middleware.RequireRole(
			"state_admin",
			
			"district_admin",
			"district_coach",

			"academy_admin",
			"academy_coach",
		))

		r.Get("/", GetDisciplinesHandler)
	})
}

