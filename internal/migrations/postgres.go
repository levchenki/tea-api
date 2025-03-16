package migrations

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/levchenki/tea-api/internal/config"
	"github.com/levchenki/tea-api/internal/storage"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunPostgresMigrations(cfg *config.Config) {
	migrationsPath := "file://migrations/postgres"

	databaseUrl := storage.GetPostgresUrl(cfg)

	m, err := migrate.New(migrationsPath, databaseUrl)
	if err != nil {
		log.Fatalf("Could not create postgres instance: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not run migrations: %v", err)
	}
	log.Println("Migrations completed")
}
