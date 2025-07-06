package teaSchemas

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/schemas"
)

type ResponseModel struct {
	Id          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	ServePrice  float64      `json:"servePrice"`
	UnitPrice   float64      `json:"unitPrice"`
	Description *string      `json:"description,omitempty"`
	CategoryId  uuid.UUID    `json:"categoryId"`
	Tags        []entity.Tag `json:"tags,omitempty"`
	IsHidden    bool         `json:"isHidden,omitempty"`
}

func NewTeaResponseModel(tea *entity.Tea) *ResponseModel {
	r := &ResponseModel{
		Id:         tea.Id,
		Name:       tea.Name,
		ServePrice: tea.ServePrice,
		UnitPrice:  tea.UnitPrice,
		CategoryId: tea.CategoryId,
	}
	if tea.Description != "" {
		r.Description = &tea.Description
	}

	if tea.Tags != nil || len(tea.Tags) > 0 {
		r.Tags = tea.Tags
	}
	if tea.IsHidden {
		r.IsHidden = tea.IsHidden
	}
	return r
}

type WithRatingResponseModel struct {
	ResponseModel
	Rating        float64 `json:"rating,omitempty"`
	AverageRating float64 `json:"averageRating,omitempty"`
	Note          string  `json:"note,omitempty" example:"This is a note"`
	IsFavourite   bool    `json:"isFavourite,omitempty"`
}

func NewTeaWithRatingResponseModel(tea *entity.TeaWithRating) *WithRatingResponseModel {
	t := &WithRatingResponseModel{
		ResponseModel: ResponseModel{
			Id:         tea.Id,
			Name:       tea.Name,
			ServePrice: tea.ServePrice,
			UnitPrice:  tea.UnitPrice,
			CategoryId: tea.CategoryId,
		},
	}
	if tea.Description != "" {
		t.Description = &tea.Description
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
	if tea.IsHidden {
		t.IsHidden = tea.IsHidden
	}
	if tea.IsFavourite {
		t.IsFavourite = tea.IsFavourite
	}
	return t
}

type MinMaxPricesResponseModel struct {
	MinServePrice float64 `json:"minServePrice"`
	MaxServePrice float64 `json:"maxServePrice"`
}

func NewMinMaxPrices(min, max float64) *MinMaxPricesResponseModel {
	return &MinMaxPricesResponseModel{
		MinServePrice: min,
		MaxServePrice: max,
	}
}

type TeaPricesPaginatedResult[T any] struct {
	schemas.PaginatedResult[T]
	MinMaxPricesResponseModel
}

func NewTeaPricesPaginatedResult[T any](teaResponse []T, total uint64, min, max float64) *TeaPricesPaginatedResult[T] {
	return &TeaPricesPaginatedResult[T]{
		PaginatedResult: schemas.PaginatedResult[T]{
			Items: teaResponse,
			Total: total,
		},
		MinMaxPricesResponseModel: MinMaxPricesResponseModel{
			MinServePrice: min,
			MaxServePrice: max,
		},
	}
}
