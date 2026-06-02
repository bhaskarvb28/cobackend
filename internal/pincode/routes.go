package pincode

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router) {
	r.Route("/pincodes", func(r chi.Router) {
		r.Get("/", GetPincodesHandler)
	})
}


