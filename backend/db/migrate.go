package db

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(connStr string) error {
	dbConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*dbConfig)
	defer db.Close()

	goose.SetBaseFS(MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
