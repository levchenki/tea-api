package entity

import (
	"github.com/google/uuid"
	"time"
)

type Tea struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	ServePrice  float64   `db:"serve_price" json:"servePrice"`
	WeightPrice float64   `db:"weight_price" json:"weightPrice"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	IsDeleted   bool      `db:"is_deleted" json:"isDeleted"`
	CategoryId  uuid.UUID `db:"category_id" json:"categoryId"`

	Tags []Tag `json:"tags,omitempty"`
}

type TeaWithRating struct {
	Tea
	Rating        float64 `db:"rating,omitempty"`
	Note          string  `db:"note,omitempty"`
	AverageRating float64 `db:"average_rating, omitempty"`
}
