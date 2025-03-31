package tagSchemas

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
)

type ResponseModel struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}

func NewResponseModel(tag *entity.Tag) *ResponseModel {
	t := &ResponseModel{
		Id:    tag.Id,
		Name:  tag.Name,
		Color: tag.Color,
	}
	return t
}
