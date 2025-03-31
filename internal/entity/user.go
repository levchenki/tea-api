package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id         uuid.UUID `db:"id"`
	TelegramId uint64    `db:"telegram_id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name,omitempty"`
	Username   string    `db:"username,omitempty"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IsAdmin    bool      `db:"is_admin"`
}

func NewEmptyUser(telegramId uint64, firstName, lastName, username string) *User {
	u := &User{
		TelegramId: telegramId,
		FirstName:  firstName,
		LastName:   lastName,
		Username:   username,
		IsAdmin:    false,
	}
	return u
}
