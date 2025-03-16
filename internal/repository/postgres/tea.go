package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/schemas"
	"log"
	"strings"
)

type TeaRepository struct {
	db *sqlx.DB
}

func NewTeaRepository(db *sqlx.DB) *TeaRepository {
	return &TeaRepository{
		db: db,
	}
}

func (r *TeaRepository) GetById(id uuid.UUID) (*entity.Tea, error) {
	tea := entity.Tea{}
	err := r.db.Get(&tea, `
		select teas.* 
		from teas
		where teas.id = $1
		limit 1`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tea, nil
}

func (r *TeaRepository) GetAll(filters *schemas.TeaFilters) ([]entity.Tea, error) {
	teas := make([]entity.Tea, 0)

	namedQuery, args, err := r.prepareFilteredQuery(filters)
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("Built named query: %s", namedQuery))

	err = r.db.Select(&teas, namedQuery, args...)
	if err != nil {
		return nil, err
	}

	return teas, nil
}

func (r *TeaRepository) prepareFilteredQuery(filters *schemas.TeaFilters) (string, []interface{}, error) {
	filterStatements := make([]string, 0, 10)
	query := "select teas.* from teas"

	if filters.CategoryId != uuid.Nil {
		categoryStmt := "category_id = :category_id"
		filterStatements = append(filterStatements, categoryStmt)
	}

	if filters.Name != "" {
		nameStmt := "name like %:name%"
		filterStatements = append(filterStatements, nameStmt)
	}

	if len(filters.Tags) > 0 {
		query += " join teas_tags tt on teas.id = tt.tea_id join tags on tt.tag_id = tags.id"

		tagStmt := "tags.id in (:tags)"
		filterStatements = append(filterStatements, tagStmt)
	}

	if filters.MinPrice != 0 && filters.MaxPrice != 0 {
		priceStmt := "price between :min_price and :max_price"
		filterStatements = append(filterStatements, priceStmt)
	}

	if len(filterStatements) > 0 {
		query += fmt.Sprintf(" where %s", strings.Join(filterStatements, " and "))
	}

	if filters.SortBy != "" {
		a := "asc"
		if filters.IsAsc {
			a = "asc"
		} else {
			a = "desc"
		}
		orderBy := fmt.Sprintf(" order by %s %s", filters.SortBy, a)

		query += orderBy
	}

	query += " limit :limit offset :offset"

	query, args, err := sqlx.Named(query, filters)
	if err != nil {
		return "", nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}

	query = r.db.Rebind(query)
	return query, args, nil
}

func (r *TeaRepository) Create(inputTea *schemas.TeaRequestModel) (*entity.Tea, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	tea := &entity.Tea{
		Name:        inputTea.Name,
		Price:       inputTea.Price,
		Description: inputTea.Description,
		CategoryId:  inputTea.CategoryId,
	}
	createdTea, err := r.insertTea(tea, tx)

	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	if len(inputTea.TagIds) != 0 {
		err := r.insertTags(&createdTea.Id, inputTea.TagIds, tx)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return createdTea, nil
}

func (r *TeaRepository) insertTea(inputTea *entity.Tea, tx *sqlx.Tx) (*entity.Tea, error) {
	createdTea := &entity.Tea{}
	rows, err := tx.NamedQuery(`
		insert into teas (name, price, description, category_id)
		values (:name, :price, :description, :category_id)
		returning teas.*`, inputTea)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		err := rows.StructScan(&createdTea)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()
	return createdTea, nil
}

func (r *TeaRepository) insertTags(teaId *uuid.UUID, tagIds []uuid.UUID, tx *sqlx.Tx) error {
	teaTags := make([]map[string]interface{}, 0, len(tagIds))
	for _, tagId := range tagIds {
		teaTags = append(teaTags, map[string]interface{}{
			"tea_id": teaId.String(),
			"tag_id": tagId.String(),
		})
	}
	_, err := tx.NamedExec(`
		insert into teas_tags (tea_id, tag_id)
		values (:tea_id, :tag_id)`, teaTags)
	if err != nil {
		return err
	}
	return nil
}

func (r *TeaRepository) deleteTags(teaId *uuid.UUID, tagIds []uuid.UUID, tx *sqlx.Tx) error {
	tagIdsStr := make([]string, len(tagIds))
	for i, tagId := range tagIds {
		tagIdsStr[i] = tagId.String()
	}

	filterParams := map[string]interface{}{
		"tea_id":  teaId.String(),
		"tag_ids": tagIdsStr,
	}

	query, args, err := sqlx.Named(
		"delete from teas_tags where tea_id=:tea_id and tag_id in (:tag_ids)",
		filterParams)

	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = tx.Rebind(query)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *TeaRepository) Update(id *uuid.UUID, inputTea *schemas.TeaRequestModel, tagsToInsert, tagsToDelete []uuid.UUID) (*entity.Tea, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	updatedTea, err := r.updateTea(id, inputTea, tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	if len(tagsToInsert) != 0 {
		err = r.insertTags(id, tagsToInsert, tx)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}
	}

	if len(tagsToDelete) != 0 {
		err = r.deleteTags(id, tagsToDelete, tx)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return nil, errRollback
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return updatedTea, nil
}

func (r *TeaRepository) updateTea(id *uuid.UUID, inputTea *schemas.TeaRequestModel, tx *sqlx.Tx) (*entity.Tea, error) {
	tea := &entity.Tea{
		Id:          *id,
		Name:        inputTea.Name,
		Price:       inputTea.Price,
		Description: inputTea.Description,
		CategoryId:  inputTea.CategoryId,
	}

	rows, err := tx.NamedQuery(`
		update teas set
			name=:name,
			price=:price,
			description=:description,
			updated_at=now(),
			category_id=:category_id
			where id = :id
		returning teas.*
		`, tea)

	if err != nil {
		return nil, err
	}

	updatedTea := &entity.Tea{}
	if rows.Next() {
		err := rows.StructScan(&updatedTea)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()
	return updatedTea, nil
}

func (r *TeaRepository) Delete(id uuid.UUID) (bool, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("delete from teas_tags where tea_id = $1", id)

	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return false, errRollback
		}
		return false, err
	}

	_, err = tx.Exec("delete from teas where id = $1", id)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return false, errRollback
		}
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}
