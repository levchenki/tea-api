package entity

import "github.com/google/uuid"

type Tag struct {
	Id    uuid.UUID `db:"id" json:"id"`
	Name  string    `db:"name" json:"name"`
	Color string    `db:"color" json:"color"`
}
