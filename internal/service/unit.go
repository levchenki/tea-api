package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
)

type UnitRepository interface {
	GetAll() ([]entity.Unit, error)
	Create(unit *entity.Unit) (*entity.Unit, error)
	Delete(id uuid.UUID) error
	Update(category *entity.Unit) (*entity.Unit, error)
	Exists(id uuid.UUID) (bool, error)
}

type UnitService struct {
	unitRepository UnitRepository
}

func NewUnitService(unitRepository UnitRepository) *UnitService {
	return &UnitService{
		unitRepository: unitRepository,
	}
}

func (s *UnitService) GetAll() ([]entity.Unit, error) {
	units, err := s.unitRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return units, err
}

func (s *UnitService) Create(unit *entity.Unit) (*entity.Unit, error) {
	createdUnit, err := s.unitRepository.Create(unit)
	if err != nil {
		return nil, err
	}

	return createdUnit, nil
}

func (s *UnitService) Update(id uuid.UUID, unit *entity.Unit) (*entity.Unit, error) {
	exists, err := s.unitRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = fmt.Errorf("unit with id %s is not found", id.String())
		return nil, errx.NewNotFoundError(err)
	}

	unit.Id = id
	updatedUnit, err := s.unitRepository.Update(unit)
	if err != nil {
		return nil, err
	}

	return updatedUnit, nil
}

func (s *UnitService) Delete(id uuid.UUID) error {
	exists, err := s.unitRepository.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		err = fmt.Errorf("unit with id %s is not found", id.String())
		return errx.NewNotFoundError(err)
	}

	err = s.unitRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
