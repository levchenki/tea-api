package v1

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	_ "github.com/levchenki/tea-api/docs"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
	"net/http"
)

type TeaService interface {
	GetTeaById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error)
	GetAllTeas(filters *teaSchemas.Filters) ([]entity.TeaWithRating, error)
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

// GetTeaById godoc
//
//	@Summary	Return tea by ID
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Tea ID"
//	@Success	200	{object}	teaSchemas.WithRatingResponseModel
//	@Failure	400	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas/{id} [get]
func (c *TeaController) GetTeaById(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	user := r.Context().Value("user")
	userClaims, ok := user.(*schemas.UserClaims)
	var userId uuid.UUID
	if ok {
		userId = userClaims.Id
	} else {
		userId = uuid.Nil
	}

	teaById, err := c.teaService.GetTeaById(id, userId)

	if err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := teaSchemas.NewTeaWithRatingResponseModel(teaById)
	render.JSON(w, r, response)
	return
}

// GetAllTeas godoc
//
//	@Summary	Get all teas
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int						false	"Page number"
//	@Param		limit		query		int						false	"Page size"
//	@Param		categoryId	query		string					false	"Category ID"
//	@Param		name		query		string					false	"Tea name"
//	@Param		tags[]		query		[]string				false	"Tags"
//	@Param		isAsc		query		bool					false	"Sort order"
//	@Param		sortBy		query		teaSchemas.SortByFilter	false	"Sort by field (name, price)"
//	@Param		price[]		query		[]float64				false	"Price range"
//	@Param		isDeleted	query		bool					false	"Is deleted"
//	@Success	200			{array}		teaSchemas.WithRatingResponseModel
//	@Failure	400			{object}	errx.AppError
//	@Failure	500			{object}	errx.AppError
//	@Router		/api/v1/teas [get]
func (c *TeaController) GetAllTeas(w http.ResponseWriter, r *http.Request) {
	filters := &teaSchemas.Filters{}

	if err := filters.Validate(r); err != nil {
		errorResponse := errx.NewBadRequestError(err)
		handleError(w, r, errorResponse)
		return
	}

	user := r.Context().Value("user")
	userClaims, ok := user.(*schemas.UserClaims)
	if ok {
		filters.UserId = userClaims.Id
	}

	teas, err := c.teaService.GetAllTeas(filters)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response := make([]*teaSchemas.WithRatingResponseModel, len(teas))
	for i := range teas {
		response[i] = teaSchemas.NewTeaWithRatingResponseModel(&teas[i])
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// CreateTea godoc
//
//	@Summary	Create tea
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		tea	body		teaSchemas.RequestModel	true	"Tea"
//	@Success	201	{object}	teaSchemas.ResponseModel
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas [post]
//	@Security	BearerAuth
func (c *TeaController) CreateTea(w http.ResponseWriter, r *http.Request) {
	teaRequest := &teaSchemas.RequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, errResponse)
		return
	}

	tea, err := c.teaService.CreateTea(teaRequest)
	if err != nil {
		handleError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, teaSchemas.NewTeaResponseModel(tea))
}

// DeleteTea godoc
//
//	@Summary	Delete tea
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Tea ID"
//	@Success	200	{object}	bool
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas/{id} [delete]
//	@Security	BearerAuth
func (c *TeaController) DeleteTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	err = c.teaService.DeleteTea(id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.JSON(w, r, true)
}

// UpdateTea godoc
//
//	@Summary	Update tea
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string					true	"Tea ID"
//	@Param		tea	body		teaSchemas.RequestModel	true	"Tea"
//	@Success	200	{object}	bool
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas/{id} [put]
//	@Security	BearerAuth
func (c *TeaController) UpdateTea(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errResponse)
		return
	}

	teaRequest := &teaSchemas.RequestModel{}
	if err := render.Bind(r, teaRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, errResponse)
		return
	}

	tea, err := c.teaService.UpdateTea(id, teaRequest)
	if err != nil {
		handleError(w, r, err)
		return
	}

	if tea == nil {
		errResponse := errx.NewNotFoundError(fmt.Errorf("tea with id %s is not found", id.String()))
		handleError(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tea)
}

// Evaluate godoc
//
//	@Summary	Evaluate tea
//	@Tags		Tea
//	@Accept		json
//	@Produce	json
//	@Param		id			path		string					true	"Tea ID"
//	@Param		evaluation	body		teaSchemas.Evaluation	true	"Evaluation"
//	@Success	200			{object}	teaSchemas.WithRatingResponseModel
//	@Failure	400			{object}	errx.AppError
//	@Failure	401			{object}	errx.AppError
//	@Failure	404			{object}	errx.AppError
//	@Failure	500			{object}	errx.AppError
//	@Router		/api/v1/teas/{id}/evaluate [post]
//	@Security	BearerAuth
func (c *TeaController) Evaluate(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errorResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, errorResponse)
		return
	}

	userClaims := r.Context().Value("user").(*schemas.UserClaims)

	evaluation := &teaSchemas.Evaluation{}
	if err = render.Bind(r, evaluation); err != nil {
		errorResponse := errx.NewBadRequestError(err)
		handleError(w, r, errorResponse)
		return
	}

	evaluatedTea, err := c.teaService.Evaluate(id, userClaims.Id, evaluation)
	if err != nil {
		handleError(w, r, err)
		return
	}

	if evaluatedTea == nil {
		errResponse := errx.NewNotFoundError(fmt.Errorf("tea with id %s is not found", id.String()))
		handleError(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, teaSchemas.NewTeaWithRatingResponseModel(evaluatedTea))
}
