package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/levchenki/tea-api/internal/config"
	"github.com/levchenki/tea-api/internal/controller"
	"github.com/levchenki/tea-api/internal/migrations"
	"github.com/levchenki/tea-api/internal/repository/postgres"
	"github.com/levchenki/tea-api/internal/service"
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

	teaRepository := postgres.NewTeaRepository(db)
	tagRepository := postgres.NewTagRepository(db)
	userRepository := postgres.NewUserRepository(db)
	categoryRepository := postgres.NewCategoryRepository(db)

	teaService := service.NewTeaService(teaRepository, tagRepository)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository, teaRepository)

	teaController := controller.NewTeaController(teaService)
	categoryController := controller.NewCategoryController(categoryService)

	authController := controller.NewAuthController(
		cfg.JWTSecretKey,
		cfg.BotToken,
		userService,
	)

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

	r.Post("/auth", authController.Auth)

	r.Route("/teas", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authController.AuthMiddleware)
			r.Get("/", teaController.GetAllTeas)
			r.Get("/{id}", teaController.GetTeaById)
			r.Post("/", teaController.CreateTea)
			r.Delete("/{id}", teaController.DeleteTea)
			r.Put("/{id}", teaController.UpdateTea)
			r.Post("/{id}/evaluate", teaController.Evaluate)
		})
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/{id}", categoryController.GetCategoryById)
		r.Get("/", categoryController.GetAllCategories)

		r.Group(func(r chi.Router) {
			r.Use(authController.AuthMiddleware)
			r.Post("/", categoryController.CreateCategory)
			r.Delete("/{id}", categoryController.DeleteCategory)
			r.Put("/{id}", categoryController.UpdateCategory)
		})
	})

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), r)
}
