package controller

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
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
}

func NewTagController(tagService TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (c *TagController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.tagService.GetAll()
	if err != nil {
		errResponse := errx.ErrorInternalServer(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
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

func (c *TagController) CreateTag(w http.ResponseWriter, r *http.Request) {
	tagRequest := &tagSchemas.RequestModel{}
	if err := render.Bind(r, tagRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	tag := &entity.Tag{
		Name:  tagRequest.Name,
		Color: tagRequest.Color,
	}

	createdTag, err := c.tagService.Create(tag)
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
	response := tagSchemas.NewResponseModel(createdTag)
	render.JSON(w, r, response)
}

func (c *TagController) UpdateTag(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	tagRequest := &tagSchemas.RequestModel{}
	if err = render.Bind(r, tagRequest); err != nil {
		errResponse := errx.ErrorBadRequest(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	tag := &entity.Tag{
		Name:  tagRequest.Name,
		Color: tagRequest.Color,
	}

	updatedTag, err := c.tagService.Update(id, tag)

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
	response := tagSchemas.NewResponseModel(updatedTag)
	render.JSON(w, r, response)
}

func (c *TagController) DeleteTag(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid id"))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	err = c.tagService.Delete(id)
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
