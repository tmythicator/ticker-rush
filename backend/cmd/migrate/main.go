// Package main runs database migrations.
package main

import (
	"log"

	"github.com/tmythicator/ticker-rush/server/db"
	"github.com/tmythicator/ticker-rush/server/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := cfg.DatabaseURL()

	log.Printf("Connecting to database...")
	if err := db.Migrate(connStr); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
