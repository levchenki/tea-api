package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/levchenki/tea-api/internal/entity"
)

type UnitRepository struct {
	db *sqlx.DB
}

func NewUnitRepository(db *sqlx.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

func (r *UnitRepository) GetById(uuid2 uuid.UUID) (*entity.Unit, error) {
	unit := &entity.Unit{}
	err := r.db.Get(unit, `
		select
		    id,
		    is_apiece,
		    weight_unit,
		    value 
		from units where id = $1`, uuid2)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (r *UnitRepository) GetAll() ([]entity.Unit, error) {
	units := make([]entity.Unit, 0)
	err := r.db.Select(&units,
		"select id, is_apiece, weight_unit, value from units order by is_apiece, weight_unit")
	if err != nil {
		return units, err
	}
	return units, nil
}

func (r *UnitRepository) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "select exists(select 1 from units where id = $1)", id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UnitRepository) Create(unit *entity.Unit) (*entity.Unit, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
        insert into units (is_apiece, weight_unit, value)
        values (:is_apiece, :weight_unit, :value)
        returning units.id, units.is_apiece, units.weight_unit, units.value
    `, unit)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	createdUnit := &entity.Unit{}
	if rows.Next() {
		err := rows.StructScan(createdUnit)
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

	return createdUnit, nil
}

func (r *UnitRepository) Update(unit *entity.Unit) (*entity.Unit, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
        update units
        set is_apiece = :is_apiece,
            weight_unit = :weight_unit,
            value = :value
        where id = :id
        returning units.id, units.is_apiece, units.weight_unit, units.value
    `, unit)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	updatedUnit := &entity.Unit{}
	if rows.Next() {
		err := rows.StructScan(updatedUnit)
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

	return updatedUnit, nil
}

func (r *UnitRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from units where id = $1", id)
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
