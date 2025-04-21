package app

import (
	"fmt"
	"github.com/levchenki/tea-api/internal/api"
	"github.com/levchenki/tea-api/internal/config"
	"github.com/levchenki/tea-api/internal/logx"
	"github.com/levchenki/tea-api/internal/logx/slogx"
	"github.com/levchenki/tea-api/internal/migrations"
	"github.com/levchenki/tea-api/internal/storage"
	"net/http"
	"os"
	"os/signal"
)

func Run() {
	cfg := config.Setup()
	var log logx.AppLogger = slogx.Setup(cfg.Environment)

	log.Info("Running migrations...")
	err := migrations.RunPostgresMigrations(cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("Migrations completed successfully")

	log.Info("Connecting to database...")
	db, err := storage.NewPostgresConnection(cfg)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("Connected to database successfully")

	r := api.NewRouter(cfg, db, log)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
	}

	log.Info("Starting server...")
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	log.Info(fmt.Sprintf("Server started on port %s", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutting down server...")
	err = server.Close()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Server shut down successfully")
	log.Info("Closing database connection...")
	err = db.Close()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Database connection closed successfully")
}
