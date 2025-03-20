package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetId(telegramId int) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.Get(&id,
		"select users.id from users where telegram_id = $1 limit 1", &telegramId)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
