package unitSchemas

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
)

type ResponseModel struct {
	Id         uuid.UUID `json:"id"`
	IsApiece   bool      `json:"isApiece"`
	WeightUnit string    `json:"weightUnit"`
	Value      int64     `json:"value"`
}

func NewResponseModel(unit *entity.Unit) *ResponseModel {
	t := &ResponseModel{
		Id:         unit.Id,
		IsApiece:   unit.IsApiece,
		WeightUnit: unit.WeightUnit.String(),
		Value:      unit.Value,
	}
	return t
}
