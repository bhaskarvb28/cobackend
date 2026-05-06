package main

import (
	"cobackend/internal/auth"
	"cobackend/internal/db"
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
	auth.RegisterRoutes(r)
	log.Println("Server running on :" + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}