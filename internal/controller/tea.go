package controller

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
	"net/http"
)

type TeaService interface {
	GetTeaById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error)
	GetAllTeas(filters *teaSchemas.Filters, userId uuid.UUID) ([]entity.TeaWithRating, error)
	CreateTea(tea *teaSchemas.RequestModel) (*entity.Tea, error)
	DeleteTea(id uuid.UUID) error
	UpdateTea(id uuid.UUID, tea *teaSchemas.RequestModel) (*entity.Tea, error)
	Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) (*entity.TeaWithRating, error)
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
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	userClaims := r.Context().Value("user").(*schemas.UserClaims)
	teaById, err := c.teaService.GetTeaById(id, userClaims.Id)

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
	response := teaSchemas.NewTeaWithRatingResponseModel(teaById)
	render.JSON(w, r, response)
}

func (c *TeaController) GetAllTeas(w http.ResponseWriter, r *http.Request) {
	filters := &teaSchemas.Filters{}

	if err := filters.Validate(r); err != nil {
		errorResponse := errx.ErrorBadRequest(err)
		render.Status(r, errorResponse.HTTPStatusCode)
		render.JSON(w, r, errorResponse)
		return
	}

	userClaims := r.Context().Value("user").(*schemas.UserClaims)

	teas, err := c.teaService.GetAllTeas(filters, userClaims.Id)
	if err != nil {
		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	response := make([]*teaSchemas.WithRatingResponseModel, len(teas))
	for i := range teas {
		response[i] = teaSchemas.NewTeaWithRatingResponseModel(&teas[i])
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func (c *TeaController) CreateTea(w http.ResponseWriter, r *http.Request) {
	teaRequest := &teaSchemas.RequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
	}

	tea, err := c.teaService.CreateTea(teaRequest)
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
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, teaSchemas.NewTeaResponseModel(tea))
}

func (c *TeaController) DeleteTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	err = c.teaService.DeleteTea(id)
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

	render.JSON(w, r, true)
}

func (c *TeaController) UpdateTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	teaRequest := &teaSchemas.RequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	tea, err := c.teaService.UpdateTea(id, teaRequest)
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

	if tea == nil {
		errResponse := errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tea)
}

func (c *TeaController) Evaluate(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	userClaims := r.Context().Value("user").(*schemas.UserClaims)

	evaluation := &teaSchemas.Evaluation{}
	if err := render.Bind(r, evaluation); err != nil {
		errorResponse := errx.ErrorBadRequest(err)
		render.Status(r, errorResponse.HTTPStatusCode)
		render.JSON(w, r, errorResponse)
		return
	}

	evaluatedTea, err := c.teaService.Evaluate(id, userClaims.Id, evaluation)
	if err != nil {
		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	if evaluatedTea == nil {
		errResponse := errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, teaSchemas.NewTeaWithRatingResponseModel(evaluatedTea))
}
