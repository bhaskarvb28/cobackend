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

	var roleName string
	err = db.QueryRow(context.Background(), `
		SELECT r.role_name 
		FROM profiles p 
		JOIN roles r ON p.role_id = r.role_id 
		WHERE p.email = 'john@gmail.com'
	`).Scan(&roleName)
	
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User 'john@gmail.com' has role: %s\n", roleName)
}
