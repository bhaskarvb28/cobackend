package main

import (
	"cobackend/internal/db"
	"log"
	"net/http"

	"fmt"

	"github.com/joho/godotenv"

	"cobackend/internal/auth"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("DB Connection Failed:", err)
	}

	fmt.Println(dbConn)

	mux := http.NewServeMux()

	auth.RegisterRoutes(mux)

	http.ListenAndServe(":8080", mux)
}