package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id         uuid.UUID `db:"id"`
	TelegramId uint64    `db:"telegram_id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func NewEmptyUser(telegramId uint64, firstName, lastName, username string) *User {
	return &User{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		Username:   username,
	}
}
