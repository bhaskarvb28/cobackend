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

	// Update role names to match code expectations (lowercase with underscores)
	_, err = db.Exec(context.Background(), "UPDATE roles SET role_name = 'super_admin' WHERE role_name = 'Super Admin'")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(context.Background(), "UPDATE roles SET role_name = 'state_admin' WHERE role_name = 'State Admin'")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(context.Background(), "UPDATE roles SET role_name = 'district_admin' WHERE role_name = 'District Admin'")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Role names updated successfully to match code expectations.")
}
