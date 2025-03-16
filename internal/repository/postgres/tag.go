package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
)

type TagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) GetAll() ([]entity.Tag, error) {
	tags := make([]entity.Tag, 0)
	err := r.db.Select(&tags, "select * from tags")
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetByTeaId(teaId uuid.UUID) ([]entity.Tag, error) {
	tags := make([]entity.Tag, 0)
	err := r.db.Select(&tags, `
		select tags.*
		from tags
		join teas_tags tt on tags.id = tt.tag_id
		where tt.tea_id = $1`, teaId)

	if err != nil {
		return nil, err
	}
	return tags, nil
}
