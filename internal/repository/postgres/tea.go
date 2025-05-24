package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
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

func (r *TeaRepository) GetById(id uuid.UUID) (*entity.TeaWithRating, error) {
	tea := entity.TeaWithRating{}
	query := `
		select
			t.id,
			name,
			serve_price,
			weight_price,
			coalesce(description, '') as description,
			t.created_at,
			t.updated_at,
			is_deleted,
			category_id,
			coalesce((select avg(rating) from evaluations where tea_id = t.id), 0) as average_rating
		from teas t
	 		left join evaluations on t.id = evaluations.tea_id
		where t.id = $1 
		limit 1`
	err := r.db.Get(&tea, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tea, nil
}
func (r *TeaRepository) GetByIdWithUser(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error) {
	tea := entity.TeaWithRating{}
	query := `
		select
			t.id,
			name,
			serve_price,
			weight_price,
			coalesce(description, '') as description,
			t.created_at,
			t.updated_at,
			is_deleted,
			category_id,
			coalesce(rating, 0) as rating,
			coalesce(note, '') as note,
			coalesce((select avg(rating) from evaluations where tea_id = t.id), 0) as average_rating
		from teas t
	 		left join evaluations on t.id = evaluations.tea_id and user_id = $1
		where t.id = $2 
		limit 1`
	err := r.db.Get(&tea, query, userId, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tea, nil
}

func (r *TeaRepository) GetAll(filters *teaSchemas.Filters) ([]entity.TeaWithRating, uint64, error) {
	teas := make([]entity.TeaWithRating, 0)

	getAllQuery, getAllArgs, err := r.prepareSelectAllQuery(filters)
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Select(&teas, getAllQuery, getAllArgs...)
	if err != nil {
		return nil, 0, err
	}

	countQuery, countArgs, err := r.prepareCountQuery(filters)
	if err != nil {
		return nil, 0, err
	}
	var totalCount uint64
	err = r.db.Get(&totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return teas, totalCount, nil
}

func (r *TeaRepository) prepareCountQuery(filters *teaSchemas.Filters) (string, []interface{}, error) {
	countQuery := "select count(*) from teas t"
	_, whereClause := r.selectAllWhereClause(countQuery, filters)
	countQuery += whereClause

	countQuery, args, err := r.bindParams(countQuery, filters)
	if err != nil {
		return "", nil, err
	}
	return countQuery, args, nil
}

func (r *TeaRepository) prepareSelectAllQuery(filters *teaSchemas.Filters) (string, []interface{}, error) {
	filterStatements := make([]string, 0, 10)
	var getAllQuery string
	isNotEmptyUser := filters.UserId != uuid.Nil
	if isNotEmptyUser {
		getAllQuery = `
		select
			t.id,
			name,
			serve_price,
			weight_price,
			coalesce(description, '') as description,
			t.created_at,
			t.updated_at,
			is_deleted,
			category_id,
			coalesce(rating, 0) as rating 
		from teas t
		    left join evaluations e on t.id = e.tea_id`
		userStmt := "(e.user_id = :user_id or rating is null)"
		filterStatements = append(filterStatements, userStmt)
	} else {
		getAllQuery = `
		select 
			id,
			name,
			serve_price,
			weight_price,
			coalesce(description, '') as description,
			created_at,
			updated_at,
			is_deleted,
			category_id
		from teas t`
	}

	getAllQuery, whereClause := r.selectAllWhereClause(getAllQuery, filters)

	getAllQuery += whereClause

	if filters.SortBy != "" {
		a := "asc"
		if filters.IsAsc {
			a = "asc"
		} else {
			a = "desc"
		}
		if isNotEmptyUser || filters.SortBy != teaSchemas.Rating {
			orderBy := fmt.Sprintf(" order by %s %s", filters.SortBy.ToDbFilter(), a)
			getAllQuery += orderBy
		}
	}

	filters.Offset = filters.Limit * (filters.Page - 1)
	getAllQuery += " limit :limit offset :offset"

	getAllQuery, args, err := r.bindParams(getAllQuery, filters)
	if err != nil {
		return "", nil, err
	}

	return getAllQuery, args, nil
}

func (r *TeaRepository) selectAllWhereClause(getQuery string, filters *teaSchemas.Filters) (string, string) {
	var whereStmt string
	filterStatements := make([]string, 0, 10)
	if filters.CategoryId != uuid.Nil {
		categoryStmt := "category_id = :category_id"
		filterStatements = append(filterStatements, categoryStmt)
	}

	if filters.Name != "" {
		nameStmt := "lower(t.name) like lower(:name)"
		filters.Name = fmt.Sprintf("%%%s%%", filters.Name)
		filterStatements = append(filterStatements, nameStmt)
	}

	if len(filters.Tags) > 0 {
		getQuery += " join teas_tags tt on t.id = tt.tea_id join tags on tt.tag_id = tags.id"

		tagStmt := "tags.id in (:tags)"
		filterStatements = append(filterStatements, tagStmt)
	}

	if filters.MinServePrice != 0 && filters.MaxServePrice != 0 {
		servePriceStmt := "serve_price between :min_serve_price and :max_serve_price"
		filterStatements = append(filterStatements, servePriceStmt)
	}

	if !filters.IsDeleted {
		isDeletedStmt := "t.is_deleted is false"
		filterStatements = append(filterStatements, isDeletedStmt)
	}

	if len(filterStatements) > 0 {
		whereStmt = fmt.Sprintf(" where %s", strings.Join(filterStatements, " and "))
	}

	return getQuery, whereStmt
}

func (r *TeaRepository) bindParams(query string, args ...interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, args)
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

func (r *TeaRepository) Create(inputTea *teaSchemas.RequestModel) (*entity.Tea, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	tea := &entity.Tea{
		Name:        inputTea.Name,
		ServePrice:  inputTea.ServePrice,
		WeightPrice: inputTea.WeightPrice,
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
		err := r.insertTags(createdTea.Id, inputTea.TagIds, tx)
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
		insert into teas (name, serve_price, weight_price, description, category_id)
		values (:name, :serve_price, :weight_price, :description, :category_id)
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

func (r *TeaRepository) insertTags(teaId uuid.UUID, tagIds []uuid.UUID, tx *sqlx.Tx) error {
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

func (r *TeaRepository) deleteTags(teaId uuid.UUID, tagIds []uuid.UUID, tx *sqlx.Tx) error {
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

func (r *TeaRepository) Update(id uuid.UUID, inputTea *teaSchemas.RequestModel, tagsToInsert, tagsToDelete []uuid.UUID) (*entity.Tea, error) {
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

func (r *TeaRepository) updateTea(id uuid.UUID, inputTea *teaSchemas.RequestModel, tx *sqlx.Tx) (*entity.Tea, error) {
	tea := &entity.Tea{
		Id:          id,
		Name:        inputTea.Name,
		ServePrice:  inputTea.ServePrice,
		WeightPrice: inputTea.WeightPrice,
		Description: inputTea.Description,
		CategoryId:  inputTea.CategoryId,
		IsDeleted:   inputTea.IsDeleted,
	}

	rows, err := tx.NamedQuery(`
		update teas set
			name=:name,
			serve_price=:serve_price,
			weight_price=:weight_price,
			description=:description,
			updated_at=now(),
			category_id=:category_id,
			is_deleted=:is_deleted
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

func (r *TeaRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from teas_tags where tea_id = $1", id)

	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}

	_, err = tx.Exec("delete from teas where id = $1", id)
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

func (r *TeaRepository) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from teas where id = $1)", id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TeaRepository) ExistsByName(existedId uuid.UUID, name string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from teas where id != $1 and name = $2 )", existedId, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TeaRepository) Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	query := `
	insert into evaluations (rating, note, created_at, updated_at, tea_id, user_id)
	values ($1, $2, now(), now(), $3, $4)
	on conflict (tea_id, user_id) do update
		set rating     = excluded.rating,
			note       = excluded.note,
			updated_at = now()
	`

	_, err = tx.Exec(query, evaluation.Rating, evaluation.Note, id, userId)
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

func (r *TeaRepository) ExistsByCategoryId(categoryId uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from teas where category_id = $1)", categoryId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TeaRepository) GetMinServePrice() (float64, error) {
	var minPrice float64

	err := r.db.Get(&minPrice, "select min(serve_price) from teas")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return minPrice, nil
}

func (r *TeaRepository) GetMaxServePrice() (float64, error) {
	var maxPrice float64

	err := r.db.Get(&maxPrice, "select max(teas.serve_price) from teas")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return maxPrice, nil
}
