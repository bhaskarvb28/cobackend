package db

import (
	"context"
	"log"
	"os"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect () (*pgxpool.Pool, error) {
	var err error

	dbUserName := os.Getenv("DB_USER_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", dbUserName, dbPassword, dbHost, dbName)

	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Connected to DB Successfully")

	return DB, nil

}