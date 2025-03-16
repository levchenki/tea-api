package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetById(id uuid.UUID) (*entity.Category, error) {
	category := &entity.Category{}
	return category, nil
}

func (r *CategoryRepository) Create(category *entity.Category) (*entity.Category, error) {
	return nil, nil
}

func (r *CategoryRepository) Update(id uuid.UUID, updatedCategory *entity.Category) (*entity.Category, error) {
	return nil, nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return nil
}
