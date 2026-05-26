package districtCoach

// import (
// 	"github.com/go-chi/chi/v5"
// 	"cobackend/internal/middleware"
// )

// func RegisterRoutes(r chi.Router) {
// 	r.Route("/district-coaches", func (r chi.Router) {
// 		r.Use(middleware.AuthMiddleware)

// 		r.Group(func(r chi.Router) {
// 			r.Use(middleware.RequireRole(
// 				"state_admin",
// 			))

// 			r.Post("/invite", InviteDistrictCoachHandler)

// 			r.Get("/", GetDistrictCoachesHandler)

// 			r.Delete("/{profile_id}", DeleteDistrictCoachHandler)

// 			r.Put("/{id}", UpdateDistrictCoachHandler)
// 		})
// 	})
// }