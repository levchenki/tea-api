package migrations

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/levchenki/tea-api/internal/config"
	"github.com/levchenki/tea-api/internal/storage"
)

func RunPostgresMigrations(cfg *config.Config) error {
	migrationsPath := "file://migrations/postgres"

	databaseUrl := storage.GetPostgresUrl(cfg)

	m, err := migrate.New(migrationsPath, databaseUrl)
	if err != nil {
		return fmt.Errorf("could not create postgres instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	return nil
}
