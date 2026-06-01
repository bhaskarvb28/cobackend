package disciplines

import (
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(r chi.Router){
	r.Route("/disciplines", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/", GetDisciplinesHandler)
	})
}

