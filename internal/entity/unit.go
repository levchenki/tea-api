package entity

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)

type WeightUnit int

const (
	Gram WeightUnit = iota + 1
	Kilogram
)

func (wu *WeightUnit) String() string {
	switch *wu {
	case Gram:
		return "G"
	case Kilogram:
		return "KG"
	default:
		return "UNKNOWN"
	}
}

func (wu WeightUnit) Value() (driver.Value, error) {
	return wu.String(), nil
}

func (wu *WeightUnit) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan %T into WeightUnit", value)
	}

	switch str {
	case "G":
		*wu = Gram
	case "KG":
		*wu = Kilogram
	default:
		return fmt.Errorf("invalid WeightUnit value: %s", str)
	}

	return nil
}

type Unit struct {
	Id         uuid.UUID  `db:"id"`
	IsApiece   bool       `db:"is_apiece"`
	WeightUnit WeightUnit `db:"weight_unit"`
	Value      int64      `db:"value"`
}

func NewUnit(isApiece bool, weightUnit string, value int64) (*Unit, error) {
	var wu WeightUnit
	err := wu.Scan(weightUnit)
	if err != nil {
		return nil, fmt.Errorf("invalid weight unit: %s", weightUnit)
	}

	return &Unit{
		IsApiece:   isApiece,
		WeightUnit: wu,
		Value:      value,
	}, nil
}
