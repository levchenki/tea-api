package categorySchemas

import (
	"github.com/levchenki/tea-api/internal/entity"
)

type ResponseModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func NewResponseModel(category *entity.Category) *ResponseModel {
	r := &ResponseModel{
		Id:   category.Id.String(),
		Name: category.Name,
	}
	if category.Description != "" {
		r.Description = category.Description
	}
	return r
}
