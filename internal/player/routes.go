package player

// import (
// 	"cobackend/internal/middleware"

// 	"github.com/go-chi/chi/v5"
// )

// func RegisterRoutes(r chi.Router) {

// 	r.Route("/players", func(r chi.Router) {

// 		r.Use(middleware.AuthMiddleware)

// 		r.Group(func(r chi.Router) {
// 			r.Use(middleware.RequireRole(
// 				"academy_admin",
// 			))

// 			r.Post("/invite", InvitePlayerHandler)
// 			// r.Post("/", CreateDistrictAdminHandler)
// 			// r.Put("/{id}", UpdateDistrictAdminHandler)
// 		})
// 	})
// }


