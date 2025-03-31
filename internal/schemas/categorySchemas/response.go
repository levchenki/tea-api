package categorySchemas

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
)

type ResponseModel struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

func NewResponseModel(category *entity.Category) *ResponseModel {
	r := &ResponseModel{
		Id:   category.Id,
		Name: category.Name,
	}
	if category.Description != "" {
		r.Description = category.Description
	}
	return r
}
