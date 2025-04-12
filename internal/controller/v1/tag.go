package v1

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/logx"
	"github.com/levchenki/tea-api/internal/schemas/tagSchemas"
	"net/http"
)

type TagService interface {
	GetAll() ([]entity.Tag, error)
	Create(tag *entity.Tag) (*entity.Tag, error)
	Update(id uuid.UUID, tag *entity.Tag) (*entity.Tag, error)
	Delete(id uuid.UUID) error
}

type TagController struct {
	tagService TagService
	log        logx.AppLogger
}

func NewTagController(tagService TagService, log logx.AppLogger) *TagController {
	return &TagController{tagService: tagService, log: log}
}

// GetAllTags godoc
//
//	@Summary	Get all tags
//	@Tags		Tags
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		tagSchemas.ResponseModel
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/tags [get]
func (c *TagController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.tagService.GetAll()
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := make([]*tagSchemas.ResponseModel, 0)
	for _, tag := range tags {
		response = append(response, tagSchemas.NewResponseModel(&tag))
	}
	render.JSON(w, r, response)
	return
}

// CreateTag godoc
//
//	@Summary	Create a new tag
//	@Tags		Tags
//	@Accept		json
//	@Produce	json
//	@Param		tag	body		tagSchemas.RequestModel	true	"Tag"
//	@Success	201	{object}	tagSchemas.ResponseModel
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/tags [post]
//	@Security	BearerAuth
func (c *TagController) CreateTag(w http.ResponseWriter, r *http.Request) {
	tagRequest := &tagSchemas.RequestModel{}
	if err := render.Bind(r, tagRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, c.log, errResponse)
		return
	}

	tag := &entity.Tag{
		Name:  tagRequest.Name,
		Color: tagRequest.Color,
	}

	createdTag, err := c.tagService.Create(tag)
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	response := tagSchemas.NewResponseModel(createdTag)
	render.JSON(w, r, response)
}

// UpdateTag godoc
//
//	@Summary	Update an existing tag
//	@Tags		Tags
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string					true	"Tag ID"
//	@Param		tag	body		tagSchemas.RequestModel	true	"Tag"
//	@Success	200	{object}	tagSchemas.ResponseModel
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/tags/{id} [put]
//	@Security	BearerAuth
func (c *TagController) UpdateTag(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, c.log, errResponse)
		return
	}

	tagRequest := &tagSchemas.RequestModel{}
	if err = render.Bind(r, tagRequest); err != nil {
		errResponse := errx.NewBadRequestError(err)
		handleError(w, r, c.log, errResponse)
		return
	}

	tag := &entity.Tag{
		Name:  tagRequest.Name,
		Color: tagRequest.Color,
	}

	updatedTag, err := c.tagService.Update(id, tag)

	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusOK)
	response := tagSchemas.NewResponseModel(updatedTag)
	render.JSON(w, r, response)
}

// DeleteTag godoc
//
//	@Summary	Delete a tag
//	@Tags		Tags
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Tag ID"
//	@Success	200	{object}	bool
//	@Failure	400	{object}	errx.AppError
//	@Failure	401	{object}	errx.AppError
//	@Failure	403	{object}	errx.AppError
//	@Failure	404	{object}	errx.AppError
//	@Failure	500	{object}	errx.AppError
//	@Router		/api/v1/tags/{id} [delete]
//	@Security	BearerAuth
func (c *TagController) DeleteTag(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid id"))
		handleError(w, r, c.log, errResponse)
		return
	}

	err = c.tagService.Delete(id)
	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, true)
}
