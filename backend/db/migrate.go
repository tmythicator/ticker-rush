package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrate applies database migrations.
func Migrate(connStr string, adminUsername, adminPasswordHash string) error {
	dbConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*dbConfig)

	defer func() {
		_ = db.Close()
	}()

	goose.SetBaseFS(MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	if err := SeedAdmin(db, adminUsername, adminPasswordHash); err != nil {
		log.Printf("Warning: failed to seed admin user: %v", err)
	}

	return nil
}

// SeedAdmin seeds an admin user into the database if one does not already exist.
func SeedAdmin(db *sql.DB, username, passwordHash string) error {
	if username == "" || passwordHash == "" {
		log.Println("Skipping admin seeding (ADMIN_USERNAME or ADMIN_PASSWORD_HASH not set)")

		return nil
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing admin user: %w", err)
	}

	if exists {
		log.Printf("Admin user '%s' already exists, skipping seeding.", username)

		return nil
	}

	insertQuery := `
		INSERT INTO users (username, password_hash, first_name, last_name, website, created_at, agb_accepted_at, is_admin)
		VALUES ($1, $2, 'Admin', 'User', '', $3, $4, TRUE)
	`
	now := time.Now()
	_, err = db.Exec(insertQuery, username, passwordHash, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	log.Printf("Successfully seeded admin user: %s", username)

	return nil
}
