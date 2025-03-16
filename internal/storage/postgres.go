package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/config"
)

func NewPostgresConnection(cfg *config.Config) (*sqlx.DB, error) {
	driver := "postgres"
	url := GetPostgresUrl(cfg)

	db, err := sqlx.Connect(driver, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func GetPostgresUrl(cfg *config.Config) string {
	databaseConfig := cfg.Database
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		databaseConfig.User,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Name,
	)
}
