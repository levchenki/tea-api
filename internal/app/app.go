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

//	@title			Tea API
//	@version		1.0
//	@description	This is a Tea API for tea cafe.

//	@contact.name	Danila Levchenko

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization

// @query.collection.format multi

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
