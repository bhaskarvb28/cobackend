package main

import (
	"cobackend/internal/academy"
	"cobackend/internal/disciplines"
	// "cobackend/internal/academyAdmin"
	// "cobackend/internal/academyCoach"
	"cobackend/internal/auth"
	"cobackend/internal/db"

	// "cobackend/internal/districtCoach"
	"cobackend/internal/district"
	"cobackend/internal/invitation"
	"cobackend/internal/profile"
	"cobackend/internal/state"

	// "cobackend/internal/districtAdmin"
	// "cobackend/internal/stateAdmin"

	// "cobackend/internal/player"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	_, err := db.Connect()
	if err != nil {
		log.Fatal("DB Connection Failed:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r)
		state.RegisterRoutes(r)
		district.RegisterRoutes(r)
		invitation.RegisterRoutes(r)
		academy.RegisterRoutes(r)

		profile.RegisterRoutes(r)
		disciplines.RegisterRoute(r)

		// stateAdmin.RegisterRoutes(r)
		// districtAdmin.RegisterRoutes(r)
		// districtCoach.RegisterRoutes(r)

		// academyAdmin.RegisterRoutes(r)

		// academyCoach.RegisterRoutes(r)

		// player.RegisterRoutes(r)

	})


	
	log.Println("Server running on :" + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}