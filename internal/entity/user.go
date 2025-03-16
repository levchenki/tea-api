package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id         uuid.UUID `db:"id"`
	TelegramId uint64    `db:"telegram_id"`
	Username   string    `db:"username"`
	Phone      string    `db:"phone"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
