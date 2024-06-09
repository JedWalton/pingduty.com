package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Init() *sql.DB {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Use environment variables
	postgresURL := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal("Failed to open connection to postgres", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to open connection to postgres", err)
	}

	log.Printf("Connected to PostgreSQL")

	return db
}
