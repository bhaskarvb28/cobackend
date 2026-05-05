package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", RegisterHandler)
	// mux.HandleFunc("/auth/login", LoginHandler)
}