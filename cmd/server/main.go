package main

import (
	"cobackend/internal/academy"
	"cobackend/internal/academyAdmin"
	"cobackend/internal/auth"
	"cobackend/internal/db"
	"cobackend/internal/districtCoach"
	"cobackend/internal/districts"
	"cobackend/internal/invitations"
	"cobackend/internal/states"

	"cobackend/internal/districtadmin"
	"cobackend/internal/stateadmin"

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
		states.RegisterRoutes(r)
		districts.RegisterRoutes(r)

		stateadmin.RegisterRoutes(r)
		districtadmin.RegisterRoutes(r)
		districtCoach.RegisterRoutes(r)

		academy.RegisterRoutes(r)
		invitations.RegisterRoutes(r)
		academyAdmin.RegisterRoutes(r)
	})


	
	log.Println("Server running on :" + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}