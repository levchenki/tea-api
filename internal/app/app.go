package app

import (
	"fmt"
	"github.com/levchenki/tea-api/internal/api"
	"github.com/levchenki/tea-api/internal/config"
	"github.com/levchenki/tea-api/internal/migrations"
	"github.com/levchenki/tea-api/internal/storage"
	"log"
	"net/http"
)

func Run() {
	cfg := config.Setup()

	migrations.RunPostgresMigrations(cfg)

	db, err := storage.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := api.NewRouter(cfg, db)

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), r)
}
