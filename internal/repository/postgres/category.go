package postgres

import (
	"database/sql"
	"errors"
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
	err := r.db.Get(category,
		"select id, name, coalesce(description, '') as description from categories where id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) GetAll() ([]entity.Category, error) {
	categories := make([]entity.Category, 0)
	err := r.db.Select(&categories,
		"select id, name, coalesce(description, '') as description from categories")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Create(category *entity.Category) (*entity.Category, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
	insert into categories (name, description)
	values (:name, :description)
	returning categories.id, categories.name, coalesce(categories.description, '') as description
	`, &category)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	createdCategory := &entity.Category{}
	if rows.Next() {
		err := rows.StructScan(&createdCategory)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}
	}
	rows.Close()
	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	return createdCategory, nil
}

func (r *CategoryRepository) Update(category *entity.Category) (*entity.Category, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
		update categories
		set name=:name,
			description=:description
		where id = :id
		returning categories.*
		`, category)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	updatedCategory := &entity.Category{}
	if rows.Next() {
		err := rows.StructScan(updatedCategory)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}
	}
	rows.Close()
	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	return updatedCategory, nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from categories where id = $1", id)
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

func (r *CategoryRepository) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from categories where id = $1)", id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CategoryRepository) ExistsByName(existedId uuid.UUID, name string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from categories where id != $1 and name = $2)", existedId, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}
