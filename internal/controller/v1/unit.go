package v1

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	_ "github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/logx"
	"github.com/levchenki/tea-api/internal/schemas/unitSchemas"
	"net/http"
)

type UnitService interface {
	GetAll() ([]entity.Unit, error)
	Create(unit *entity.Unit) (*entity.Unit, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, category *entity.Unit) (*entity.Unit, error)
}

type UnitController struct {
	unitService UnitService
	log         logx.AppLogger
}

func NewUnitController(unitService UnitService, log logx.AppLogger) *UnitController {
	return &UnitController{
		unitService: unitService,
		log:         log,
	}
}

// GetAllUnits godoc
//
//	@Summary	Return all units
//	@Tags		Unit
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]unitSchemas.ResponseModel
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas/units [get]
func (c *UnitController) GetAllUnits(w http.ResponseWriter, r *http.Request) {
	units, err := c.unitService.GetAll()
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	response := make([]*unitSchemas.ResponseModel, 0)
	for _, unit := range units {
		response = append(response, unitSchemas.NewResponseModel(&unit))
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GetAllWeights godoc
//
//	@Summary	Return all weights
//	@Tags		Unit
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]string
//	@Router		/api/v1/teas/units/weights [get]
func (c *UnitController) GetAllWeights(w http.ResponseWriter, r *http.Request) {
	weightUnits := []entity.WeightUnit{entity.Gram, entity.Kilogram}
	response := make([]string, len(weightUnits))
	for i, w := range weightUnits {
		response[i] = w.String()
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// CreateUnit godoc
//
//	@Summary	Create unit
//	@Tags		Unit
//	@Accept		json
//	@Produce	json
//	@Param		unit	body		unitSchemas.RequestModel	true	"Unit"
//	@Success	201		{object}	unitSchemas.ResponseModel
//	@Failure	400		{object}	errx.AppError
//	@Failure	500		{object}	errx.AppError
//	@Router		/api/v1/teas/units [post]
func (c *UnitController) CreateUnit(w http.ResponseWriter, r *http.Request) {
	unitRequest := &unitSchemas.RequestModel{}
	if err := render.Bind(r, unitRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, c.log, errResponse)
		return
	}

	unit, err := entity.NewUnit(unitRequest.IsApiece, unitRequest.WeightUnit, unitRequest.Value)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid unit data: %w", err))
		handleError(w, r, c.log, errResponse)
		return
	}

	createdUnit, err := c.unitService.Create(unit)
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	response := unitSchemas.NewResponseModel(createdUnit)
	render.JSON(w, r, response)
}

// UpdateUnit godoc
//
//	@Summary	Update unit
//	@Tags		Unit
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"Unit ID"
//	@Param		unit	body		unitSchemas.RequestModel	true	"Unit"
//	@Success	200		{object}	unitSchemas.ResponseModel
//	@Failure	400		{object}	errx.AppError
//	@Failure	404		{object}	errx.AppError
//	@Failure	500		{object}	errx.AppError
//	@Router		/api/v1/teas/units/{id} [put]
func (c *UnitController) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, c.log, errResponse)
		return
	}

	unitRequest := &unitSchemas.RequestModel{}
	if err := render.Bind(r, unitRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, c.log, errResponse)
		return
	}

	unit, err := entity.NewUnit(unitRequest.IsApiece, unitRequest.WeightUnit, unitRequest.Value)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid unit data: %w", err))
		handleError(w, r, c.log, errResponse)
		return
	}

	updatedUnit, err := c.unitService.Update(id, unit)
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := unitSchemas.NewResponseModel(updatedUnit)
	render.JSON(w, r, response)
}

// DeleteUnit godoc
//
//	@Summary	Delete unit
//	@Tags		Unit
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Unit ID"
//	@Success	200	{object}	bool
//	@Failure	400	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/teas/units/{id} [delete]
func (c *UnitController) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, c.log, errResponse)
		return
	}

	err = c.unitService.Delete(id)
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, true)
}
