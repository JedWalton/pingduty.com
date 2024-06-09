package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// LoggerAdapter adapts a *log.Logger to meet the migrate.Logger interface.
type LoggerAdapter struct {
	*log.Logger
}

// Verbose should return true if verbose logging output is enabled.
func (l LoggerAdapter) Verbose() bool {
	// Update this to reflect whether you want verbose logging or not
	return true
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	connectionString := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to start migration: %v", err)
	}

	// Create and use the logger adapter
	logger := LoggerAdapter{Logger: log.New(log.Writer(), "MIGRATE: ", log.LstdFlags)}
	m.Log = logger

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No migration needed.")
	} else {
		log.Println("Migrations applied successfully!")
	}
}
