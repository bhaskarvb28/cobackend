package profile

import (
	"github.com/go-chi/chi/v5"

	"cobackend/internal/middleware"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/me", func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Get("/profile", GetProfileHandler)

		r.Patch("/profile", CompleteProfileHandler)
		r.Post("/profile", CompleteProfileHandler)
	})
}