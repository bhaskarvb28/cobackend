package stateAdmin

// import (
// 	"cobackend/internal/middleware"

// 	"github.com/go-chi/chi/v5"
// )

// func RegisterRoutes(r chi.Router) {

// 	r.Route("/state-admins", func(r chi.Router) {

// 		r.Use(middleware.AuthMiddleware)

// 		r.Group(func(r chi.Router) {

// 			r.Use(
// 				middleware.RequireRole("super_admin"),
// 			)

// 			r.Post("/invite", InviteStateAdminHandler)

// 			r.Get("/", GetStateAdminsHandler)
			
// 			r.Put("/{profile_id}/state", UpdateAssignedStateHandler)
// 			r.Delete("/{profile_id}", DeleteStateAdminHandler)
// 		})
// 	})
// }