// Package db provides database migrations and utilities.
package db

import "embed"

// MigrationsFS holds the embedded migration files.
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS
