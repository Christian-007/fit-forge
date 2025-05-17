package testutil

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(ctx context.Context, dbUrl string, migrationFilePath string) error {
	config, err := pgx.ParseConfig(dbUrl)
	if err != nil {
		return fmt.Errorf("error parsing dbUrl: %w", err)
	}

	sqlDB := stdlib.OpenDB(*config)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationFilePath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error executing migration up: %w", err)
	}

	return nil
}
