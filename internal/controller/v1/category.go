package v1

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas/categorySchemas"
	"net/http"
)

type CategoryService interface {
	GetById(id uuid.UUID) (*entity.Category, error)
	GetAll() ([]entity.Category, error)
	Create(category *entity.Category) (*entity.Category, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, category *entity.Category) (*entity.Category, error)
}

type CategoryController struct {
	categoryService CategoryService
}

func NewCategoryController(categoryService CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	category, err := c.categoryService.GetById(id)
	if err != nil {
		var errorResponse *errx.ErrorResponse
		if errors.As(err, &errorResponse) {
			render.Status(r, err.(*errx.ErrorResponse).HTTPStatusCode)
			render.JSON(w, r, err)
			return
		}

		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	response := categorySchemas.NewResponseModel(category)
	render.JSON(w, r, response)
}

func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categoryService.GetAll()
	if err != nil {
		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	response := make([]*categorySchemas.ResponseModel, 0)
	for _, category := range categories {
		response = append(response, categorySchemas.NewResponseModel(&category))
	}
	render.JSON(w, r, response)
}

func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	categoryRequest := &categorySchemas.RequestModel{}
	if err := render.Bind(r, categoryRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	category := &entity.Category{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}

	createdCategory, err := c.categoryService.Create(category)
	if err != nil {
		var errorResponse *errx.ErrorResponse
		if errors.As(err, &errorResponse) {
			render.Status(r, errorResponse.HTTPStatusCode)
			render.JSON(w, r, err)
			return
		}
		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusCreated)
	response := categorySchemas.NewResponseModel(createdCategory)
	render.JSON(w, r, response)
}

func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	categoryRequest := &categorySchemas.RequestModel{}
	if err := render.Bind(r, categoryRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	category := &entity.Category{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}

	updatedCategory, err := c.categoryService.Update(id, category)
	if err != nil {
		var errorResponse *errx.ErrorResponse
		if errors.As(err, &errorResponse) {
			render.Status(r, errorResponse.HTTPStatusCode)
			render.JSON(w, r, err)
			return
		}

		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	response := categorySchemas.NewResponseModel(updatedCategory)
	render.JSON(w, r, response)
}

func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	err = c.categoryService.Delete(id)
	if err != nil {
		var errorResponse *errx.ErrorResponse
		if errors.As(err, &errorResponse) {
			render.Status(r, errorResponse.HTTPStatusCode)
			render.JSON(w, r, err)
			return
		}

		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, true)
}
