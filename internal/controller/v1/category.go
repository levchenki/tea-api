package v1

import (
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

// GetCategoryById godoc
//
//	@Summary	Return category by ID
//	@Tags		Category
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Category ID"
//	@Success	200	{object}	categorySchemas.ResponseModel
//	@Failure	400	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/categories/{id} [get]
func (c *CategoryController) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	category, err := c.categoryService.GetById(id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := categorySchemas.NewResponseModel(category)
	render.JSON(w, r, response)
}

// GetAllCategories godoc
//
//	@Summary	Return all categories
//	@Tags		Category
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]categorySchemas.ResponseModel
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/categories [get]
func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categoryService.GetAll()
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := make([]*categorySchemas.ResponseModel, 0)
	for _, category := range categories {
		response = append(response, categorySchemas.NewResponseModel(&category))
	}
	render.JSON(w, r, response)
}

// CreateCategory godoc
//
//	@Summary	Create category
//	@Tags		Category
//	@Accept		json
//	@Produce	json
//	@Param		category	body		categorySchemas.RequestModel	true	"Category"
//	@Success	201			{object}	categorySchemas.ResponseModel
//	@Failure	400			{object}	errx.AppError
//	@Failure	401			{object}	errx.AppError
//	@Failure	403			{object}	errx.AppError
//	@Failure	500			{object}	errx.AppError
//	@Router		/api/v1/categories [post]
//	@Security	BearerAuth
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	categoryRequest := &categorySchemas.RequestModel{}
	if err := render.Bind(r, categoryRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, errResponse)
		return
	}

	category := &entity.Category{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}

	createdCategory, err := c.categoryService.Create(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	response := categorySchemas.NewResponseModel(createdCategory)
	render.JSON(w, r, response)
}

// UpdateCategory godoc
//
//	@Summary	Update category
//	@Tags		Category
//	@Accept		json
//	@Produce	json
//	@Param		id			path		string							true	"Category ID"
//	@Param		category	body		categorySchemas.RequestModel	true	"Category"
//	@Success	200			{object}	categorySchemas.ResponseModel
//	@Failure	400			{object}	errx.AppError
//	@Failure	401			{object}	errx.AppError
//	@Failure	403			{object}	errx.AppError
//	@Failure	404			{object}	errx.AppError
//	@Failure	500			{object}	errx.AppError
//	@Router		/api/v1/categories/{id} [put]
//	@Security	BearerAuth
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	categoryRequest := &categorySchemas.RequestModel{}
	if err := render.Bind(r, categoryRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, errResponse)
		return
	}

	category := &entity.Category{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}

	updatedCategory, err := c.categoryService.Update(id, category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := categorySchemas.NewResponseModel(updatedCategory)
	render.JSON(w, r, response)
}

// DeleteCategory godoc
//
//	@Summary	Delete category
//	@Tags		Category
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Category ID"
//	@Success	200	{object}	bool
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/categories/{id} [delete]
//	@Security	BearerAuth
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	err = c.categoryService.Delete(id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, true)
}
