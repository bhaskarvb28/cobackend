package invitation

import (
	"github.com/go-chi/chi/v5"
	
	"cobackend/internal/middleware"

)

// RegisterRoutes registers all invitation-related routes.
//
// Routes:
//
//	Public:
//		POST   /invitations/accept
//		GET    /invitations/{token}
//
//	Protected:
//		POST   /invitations/invite
//		GET    /invitations
//		GET    /invitations/{id}
//		POST   /invitations/{id}/revoke
//
// Middleware:
//
//	Protected Routes:
//		- AuthMiddleware
//		- RequireRole(
//			super_admin,
//			state_admin,
//			district_admin,
//			academy_admin,
//		  )
func RegisterRoutes(r chi.Router) {

	r.Route("/invitations", func(r chi.Router) {

		// ----------------------------------------------------------
		// Public Routes
		// ----------------------------------------------------------

		// Accept invitation using invitation token.
		r.Post(
			"/accept",
			AcceptInvitationHandler,
		)

		// Get invitation details using invitation token.
		r.Get(
			"/token/{token}",
			GetInvitationByTokenHandler,
		)

		// ----------------------------------------------------------
		// Protected Routes
		// ----------------------------------------------------------

		r.Group(func(r chi.Router) {

			r.Use(
				middleware.AuthMiddleware,
			)

			r.Use(
				middleware.RequireRole(
					"super_admin",
					"state_admin",
					"district_admin",
					"academy_admin",
				),
			)

			r.Use(
				middleware.RequireCompletedProfile,
			)

			// Get all invitations.
			r.Get(
				"/",
				GetInvitationsHandler,
			)


			// Create a new invitation.
			r.Post(
				"/invite",
				CreateInviteHandler,
			)

			r.Delete(
				"/{id}",
				DeleteInvitationHandler,
			)

			// Get invitation by ID.
			// r.Get(
			// 	"/{id}",
			// 	GetInvitationByIDHandler,
			// )

			// Revoke invitation by ID.
			
		})
	})
}