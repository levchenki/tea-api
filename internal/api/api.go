package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/levchenki/tea-api/docs"
	v1 "github.com/levchenki/tea-api/internal/api/v1"
	"github.com/levchenki/tea-api/internal/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(cfg *config.Config, db *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := v1.NewRouter(cfg, db)

	r.Get("/swagger/*", httpSwagger.Handler())

	r.Route("/api", func(r chi.Router) {

		r.Mount("/v1", v1Router)
	})

	return r
}
