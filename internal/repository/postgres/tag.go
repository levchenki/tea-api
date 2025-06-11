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
		where tt.tea_id = $1`, teaId.String())

	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) Create(tag *entity.Tag) (*entity.Tag, error) {
	createdTag := &entity.Tag{}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
	insert into tags (name, color)
	values (:name, :color)
	returning tags.*
	`, &tag)

	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	if rows.Next() {
		err := rows.StructScan(&createdTag)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	rows.Close()
	return createdTag, nil
}

func (r *TagRepository) Update(tag *entity.Tag) (*entity.Tag, error) {
	updatedTag := &entity.Tag{}
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
		update tags
		set name=:name,
			color=:color
		where id = :id
		returning tags.*
	`, tag)

	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	if rows.Next() {
		err := rows.StructScan(&updatedTag)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	rows.Close()
	return updatedTag, nil
}

func (r *TagRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from tags where id = $1`, id)
	if err != nil {

	}

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

func (r *TagRepository) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from tags where id = $1)", id)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (r *TagRepository) ExistsByName(existedId uuid.UUID, name string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from tags where id != $1 and name = $2 )",
		existedId, name)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (r *TagRepository) ExistsByTeas(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from teas_tags where teas_tags.tag_id = $1)", id)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (r *TagRepository) GetAllByTeaIds(teaIds []uuid.UUID) (map[uuid.UUID][]entity.Tag, error) {
	query := `
		select tt.tea_id,
		       t.id,
		       t.name,
		       t.color
		from tags t
		         join teas_tags tt on t.id = tt.tag_id
		where tt.tea_id in (?)
	`

	query, args, err := sqlx.In(query, teaIds)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]entity.Tag)
	for rows.Next() {
		var (
			teaId uuid.UUID
			tag   entity.Tag
		)
		err = rows.Scan(&teaId, &tag.Id, &tag.Name, &tag.Color)
		if err != nil {
			return nil, err
		}
		result[teaId] = append(result[teaId], tag)
	}

	return result, nil
}
