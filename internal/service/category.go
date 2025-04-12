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
		err := fmt.Errorf("category with id %s is not found", id.String())
		return nil, errx.NewNotFoundError(err)
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
		err := fmt.Errorf("category with name %s has already existed", category.Name)
		return nil, errx.NewBadRequestError(err)
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
		err := fmt.Errorf("category with id %s is not found", category.Id.String())
		return nil, errx.NewNotFoundError(err)
	}

	existsByName, err := s.categoryRepository.ExistsByName(id, category.Name)
	if err != nil {
		return nil, err
	}
	if existsByName == true {
		err := fmt.Errorf("category with name %s has already existed", category.Name)
		return nil, errx.NewBadRequestError(err)
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
		err := fmt.Errorf("category with id %s is not found", id.String())
		return errx.NewNotFoundError(err)
	}

	existsByCategoryId, err := s.teaRepository.ExistsByCategoryId(id)
	if err != nil {
		return err
	}
	if existsByCategoryId == true {
		err := fmt.Errorf("category with id %s has some teas", id.String())
		return errx.NewBadRequestError(err)
	}

	err = s.categoryRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
