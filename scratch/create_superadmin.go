package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	password := "password"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	superAdminRoleID := "b0b99bdc-4520-434f-9ecb-26ad96f45e95"
	email := "admin@gmail.com"

	_, err = db.Exec(context.Background(), `
		INSERT INTO profiles (id, email, password, role_id, first_name, last_name, contact_number) 
		VALUES (gen_random_uuid(), $1, $2, $3, 'Super', 'Admin', '1234567890')
		ON CONFLICT (email) DO UPDATE SET password = $2, role_id = $3
	`, email, string(hashed), superAdminRoleID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Super Admin created: %s / %s\n", email, password)
}
