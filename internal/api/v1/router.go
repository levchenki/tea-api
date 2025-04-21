package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/config"
	v1 "github.com/levchenki/tea-api/internal/controller/v1"
	"github.com/levchenki/tea-api/internal/logx"
	"github.com/levchenki/tea-api/internal/repository/postgres"
	"github.com/levchenki/tea-api/internal/service"
)

func NewRouter(cfg *config.Config, db *sqlx.DB, log logx.AppLogger) *chi.Mux {
	teaRepository := postgres.NewTeaRepository(db)
	tagRepository := postgres.NewTagRepository(db)
	userRepository := postgres.NewUserRepository(db)
	categoryRepository := postgres.NewCategoryRepository(db)

	teaService := service.NewTeaService(teaRepository, tagRepository)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository, teaRepository)
	tagService := service.NewTagService(tagRepository)

	teaControllerV1 := v1.NewTeaController(teaService, log)
	categoryControllerV1 := v1.NewCategoryController(categoryService, log)
	tagControllerV1 := v1.NewTagController(tagService, log)

	authControllerV1 := v1.NewUserController(
		cfg.JWTSecretKey,
		cfg.BotToken,
		userService,
		log,
	)

	r := chi.NewRouter()

	r.Post("/auth", authControllerV1.Auth)
	r.Route("/teas", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(authControllerV1.AuthMiddleware(false))
			r.Get("/", teaControllerV1.GetAllTeas)
			r.Get("/{id}", teaControllerV1.GetTeaById)
		})

		r.Group(func(r chi.Router) {
			r.Use(authControllerV1.AuthMiddleware(true))
			r.Post("/{id}/evaluate", teaControllerV1.Evaluate)

			r.Group(func(r chi.Router) {
				r.Use(authControllerV1.AdminMiddleware)
				r.Post("/", teaControllerV1.CreateTea)
				r.Delete("/{id}", teaControllerV1.DeleteTea)
				r.Put("/{id}", teaControllerV1.UpdateTea)
			})
		})
	})
	r.Route("/categories", func(r chi.Router) {
		r.Get("/{id}", categoryControllerV1.GetCategoryById)
		r.Get("/", categoryControllerV1.GetAllCategories)

		r.Group(func(r chi.Router) {
			r.Use(authControllerV1.AuthMiddleware(true))
			r.Use(authControllerV1.AdminMiddleware)
			r.Post("/", categoryControllerV1.CreateCategory)
			r.Delete("/{id}", categoryControllerV1.DeleteCategory)
			r.Put("/{id}", categoryControllerV1.UpdateCategory)
		})
	})
	r.Route("/tags", func(r chi.Router) {
		r.Get("/", tagControllerV1.GetAllTags)

		r.Group(func(r chi.Router) {
			r.Use(authControllerV1.AuthMiddleware(true))
			r.Use(authControllerV1.AdminMiddleware)
			r.Post("/", tagControllerV1.CreateTag)
			r.Delete("/{id}", tagControllerV1.DeleteTag)
			r.Put("/{id}", tagControllerV1.UpdateTag)
		})
	})
	return r
}
