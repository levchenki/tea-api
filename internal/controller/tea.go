package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/api"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/schemas"
	"net/http"
)

type TeaService interface {
	GetTeaById(id uuid.UUID) (*entity.Tea, error)
	GetAllTeas(filters *schemas.TeaFilters) ([]entity.Tea, error)
	CreateTea(tea *schemas.TeaRequestModel) (*entity.Tea, error)
	DeleteTea(id uuid.UUID) (bool, error)
	UpdateTea(id uuid.UUID, tea *schemas.TeaRequestModel) (*entity.Tea, error)
}

type TeaController struct {
	teaService TeaService
}

func NewTeaController(teaService TeaService) *TeaController {
	return &TeaController{teaService: teaService}
}

func (c *TeaController) GetTeaById(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := api.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	teaById, err := c.teaService.GetTeaById(id)
	if err != nil {
		errResponse := api.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	if teaById == nil {
		errResponse := api.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, schemas.NewTeaResponseModel(teaById))
}

func (c *TeaController) GetAllTeas(w http.ResponseWriter, r *http.Request) {
	filters := &schemas.TeaFilters{}
	claims := r.Context().Value("user")
	fmt.Println(claims)
	if err := filters.Validate(r); err != nil {
		errorResponse := api.ErrorBadRequest(err)
		render.Status(r, errorResponse.HTTPStatusCode)
		render.JSON(w, r, errorResponse)
		return
	}

	teas, err := c.teaService.GetAllTeas(filters)
	if err != nil {
		errResponse := api.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	response := make([]*schemas.TeaResponseModel, len(teas))
	for i := range teas {
		response[i] = schemas.NewTeaResponseModel(&teas[i])
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func (c *TeaController) CreateTea(w http.ResponseWriter, r *http.Request) {
	teaRequest := &schemas.TeaRequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := api.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}

	tea, err := c.teaService.CreateTea(teaRequest)
	if err != nil {
		errResponse := api.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, schemas.NewTeaResponseModel(tea))
}

func (c *TeaController) DeleteTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := api.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	isDeleted, err := c.teaService.DeleteTea(id)
	if err != nil {
		errResponse := api.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}

	render.JSON(w, r, isDeleted)
}

func (c *TeaController) UpdateTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := api.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	teaRequest := &schemas.TeaRequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := api.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}

	tea, err := c.teaService.UpdateTea(id, teaRequest)
	if err != nil {
		errResponse := api.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}

	if tea == nil {
		errResponse := api.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tea)
}
