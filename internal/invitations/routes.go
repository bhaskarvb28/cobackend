package invitations

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/invitations", func (r chi.Router) {

		r.Get("/{token}", GetInvitationByTokenHandler)
	})
}