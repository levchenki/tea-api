package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(user *entity.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		insert into users (telegram_id,
						   username,
						   created_at,
						   updated_at,
						   first_name,
						   last_name,
		                   is_admin)
		values (:telegram_id,
				:username,
				now(),
				now(),
				:first_name,
				:last_name,
		        :is_admin)
		`, &user)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Exists(telegramId uint64) (bool, error) {
	var exists bool
	err := r.db.Get(&exists,
		"select exists(select 1 from users where telegram_id = $1)", &telegramId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetByTelegramId(telegramId uint64) (*entity.User, error) {
	user := entity.User{}
	err := r.db.Get(&user, `
	select
		u.id,
		u.telegram_id,
		coalesce(u.first_name, '') as first_name,
		coalesce(u.last_name, '') as last_name,
		coalesce(u.username, '') as username,
		u.created_at,
		u.updated_at,
		u.is_admin
	from users u where telegram_id = $1 limit 1`, &telegramId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SaveRefreshToken(userId, refreshTokenId uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"update users set refresh_token_id = $1 where id = $2", refreshTokenId, userId,
	)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) IsRefreshTokenExists(userId, refreshTokenId uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from users where id = $1 and refresh_token_id = $2)", userId, refreshTokenId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetById(id uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.Get(user, `
	select
		u.id,
		u.telegram_id,
		coalesce(u.first_name, '') as first_name,
		coalesce(u.last_name, '') as last_name,
		coalesce(u.username, '') as username,
		u.created_at,
		u.updated_at,
		u.is_admin
	from users u where id = $1 limit 1`, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
