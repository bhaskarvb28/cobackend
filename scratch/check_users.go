package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbUserName := os.Getenv("DB_USER_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", dbUserName, dbPassword, dbHost, dbName)

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT email FROM profiles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Available Emails in Database:")
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}
		fmt.Println("-", email)
	}
}
