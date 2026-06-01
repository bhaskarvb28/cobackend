package main

import (
	"cobackend/internal/academy"
	"cobackend/internal/disciplines"
	"cobackend/internal/role"
	session "cobackend/internal/sessions"

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

	"cobackend/internal/player"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://yourfrontend.onrender.com",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},

		ExposedHeaders: []string{
			"Link",
		},

		AllowCredentials: true,

		MaxAge: 300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r)
		role.RegisterRoutes(r)
		state.RegisterRoutes(r)
		district.RegisterRoutes(r)
		invitation.RegisterRoutes(r)
		academy.RegisterRoutes(r)

		profile.RegisterRoutes(r)
		disciplines.RegisterRoute(r)

		player.RegisterRoutes(r)
		session.RegisterRoutes(r)


		// stateAdmin.RegisterRoutes(r)
		// districtAdmin.RegisterRoutes(r)
		// districtCoach.RegisterRoutes(r)

		// academyAdmin.RegisterRoutes(r)

		// academyCoach.RegisterRoutes(r)


	})


	
	log.Println("Server running on :" + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}