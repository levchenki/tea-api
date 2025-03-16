package entity

import "github.com/google/uuid"

type Category struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}
