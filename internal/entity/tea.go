package entity

import (
	"github.com/google/uuid"
	"time"
)

type Tea struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Price       float64   `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	IsDeleted   bool      `db:"is_deleted" json:"isDeleted"`
	CategoryId  uuid.UUID `db:"category_id" json:"categoryId"`

	Tags []Tag `json:"tags,omitempty"`
}
