package teaSchemas

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
)

type ResponseModel struct {
	Id          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	ServePrice  float64      `json:"serve_price"`
	WeightPrice float64      `json:"weight_price"`
	Description *string      `json:"description"`
	CategoryId  uuid.UUID    `json:"categoryId"`
	Tags        []entity.Tag `json:"tags,omitempty"`
	IsDeleted   bool         `json:"isDeleted,omitempty"`
}

func NewTeaResponseModel(tea *entity.Tea) *ResponseModel {
	r := &ResponseModel{
		Id:          tea.Id,
		Name:        tea.Name,
		ServePrice:  tea.ServePrice,
		WeightPrice: tea.WeightPrice,
		Description: &tea.Description,
		CategoryId:  tea.CategoryId,
	}
	if tea.Tags != nil || len(tea.Tags) > 0 {
		r.Tags = tea.Tags
	}
	if tea.IsDeleted {
		r.IsDeleted = tea.IsDeleted
	}
	return r
}

type WithRatingResponseModel struct {
	ResponseModel
	Rating        float64 `json:"rating,omitempty"`
	AverageRating float64 `json:"averageRating,omitempty"`
	Note          string  `json:"note,omitempty" example:"This is a note"`
}

func NewTeaWithRatingResponseModel(tea *entity.TeaWithRating) *WithRatingResponseModel {
	t := &WithRatingResponseModel{
		ResponseModel: ResponseModel{
			Id:          tea.Id,
			Name:        tea.Name,
			ServePrice:  tea.ServePrice,
			WeightPrice: tea.WeightPrice,
			CategoryId:  tea.CategoryId,
		},
	}
	if tea.Description != "" {
		t.Description = nil
	}
	if tea.Tags != nil || len(tea.Tags) > 0 {
		t.Tags = tea.Tags
	}
	if tea.Rating != 0 {
		t.Rating = tea.Rating
	}
	if tea.AverageRating != 0 {
		t.AverageRating = tea.AverageRating
	}
	if tea.Note != "" {
		t.Note = tea.Note
	}
	if tea.IsDeleted {
		t.IsDeleted = tea.IsDeleted
	}
	return t
}
