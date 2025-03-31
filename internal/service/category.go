package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
)

type CategoryRepository interface {
	GetById(id uuid.UUID) (*entity.Category, error)
	GetAll() ([]entity.Category, error)
	Create(category *entity.Category) (*entity.Category, error)
	Update(category *entity.Category) (*entity.Category, error)
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(existedId uuid.UUID, name string) (bool, error)
}

type CategoryTeaRepository interface {
	ExistsByCategoryId(categoryId uuid.UUID) (bool, error)
}

type CategoryService struct {
	categoryRepository CategoryRepository
	teaRepository      CategoryTeaRepository
}

func NewCategoryService(categoryRepository CategoryRepository, teaRepository CategoryTeaRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
		teaRepository:      teaRepository,
	}
}

func (s *CategoryService) GetById(id uuid.UUID) (*entity.Category, error) {
	category, err := s.categoryRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errx.ErrorNotFound(fmt.Errorf("category with id %s is not found", id.String()))
	}

	return category, nil
}

func (s *CategoryService) GetAll() ([]entity.Category, error) {
	categories, err := s.categoryRepository.GetAll()
	if err != nil {
		return make([]entity.Category, 0), err
	}

	return categories, nil
}

func (s *CategoryService) Create(category *entity.Category) (*entity.Category, error) {
	exists, err := s.categoryRepository.ExistsByName(uuid.Nil, category.Name)
	if err != nil {
		return nil, err
	}
	if exists == true {
		return nil, errx.ErrorBadRequest(fmt.Errorf("category with name %s has already existed", category.Name))
	}
	createdCategory, err := s.categoryRepository.Create(category)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}

func (s *CategoryService) Update(id uuid.UUID, category *entity.Category) (*entity.Category, error) {
	exists, err := s.categoryRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, errx.ErrorNotFound(fmt.Errorf("category with id %s is not found", category.Id.String()))
	}

	existsByName, err := s.categoryRepository.ExistsByName(id, category.Name)
	if err != nil {
		return nil, err
	}
	if existsByName == true {
		return nil, errx.ErrorBadRequest(fmt.Errorf("category with name %s has already existed", category.Name))
	}

	category.Id = id
	updatedCategory, err := s.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	exists, err := s.categoryRepository.Exists(id)
	if err != nil {
		return err
	}
	if exists == false {
		return errx.ErrorNotFound(fmt.Errorf("category with id %s is not found", id.String()))
	}

	existsByCategoryId, err := s.teaRepository.ExistsByCategoryId(id)
	if err != nil {
		return err
	}
	if existsByCategoryId == true {
		return errx.ErrorBadRequest(fmt.Errorf("category with id %s has some teas", id.String()))
	}

	err = s.categoryRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
