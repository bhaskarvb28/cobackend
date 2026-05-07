package states

import (
	"github.com/go-chi/chi/v5"

	"cobackend/internal/middleware"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/states", func(r chi.Router) {

		// all state routes require authentication
		r.Use(middleware.AuthMiddleware)

		r.Get("/", GetStatesHandler)
	})
}