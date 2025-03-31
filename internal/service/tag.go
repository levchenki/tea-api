package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
)

type TagRepository interface {
	GetAll() ([]entity.Tag, error)
	Create(tag *entity.Tag) (*entity.Tag, error)
	Update(tag *entity.Tag) (*entity.Tag, error)
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(existedId uuid.UUID, name string) (bool, error)
	ExistsByTeas(id uuid.UUID) (bool, error)
}

type TagService struct {
	tagRepository TagRepository
}

func NewTagService(tagRepository TagRepository) *TagService {
	return &TagService{
		tagRepository: tagRepository,
	}
}

func (s *TagService) GetAll() ([]entity.Tag, error) {
	tags, err := s.tagRepository.GetAll()
	if err != nil {
		return make([]entity.Tag, 0), err
	}
	return tags, nil
}

func (s *TagService) Create(tag *entity.Tag) (*entity.Tag, error) {
	exists, err := s.tagRepository.ExistsByName(uuid.Nil, tag.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errx.ErrorBadRequest(fmt.Errorf("tag with name %s already is exist", tag.Name))
	}

	createdTag, err := s.tagRepository.Create(tag)
	if err != nil {
		return nil, err
	}
	return createdTag, nil
}

func (s *TagService) Update(id uuid.UUID, tag *entity.Tag) (*entity.Tag, error) {

	exists, err := s.tagRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errx.ErrorNotFound(fmt.Errorf("tag with id %s is not found", id.String()))
	}

	existsByName, err := s.tagRepository.ExistsByName(id, tag.Name)
	if err != nil {
		return nil, err
	}
	if existsByName {
		return nil, errx.ErrorBadRequest(fmt.Errorf("tag with name %s already is exist", tag.Name))
	}

	tag.Id = id
	updatedTag, err := s.tagRepository.Update(tag)
	if err != nil {
		return nil, err
	}
	return updatedTag, nil
}

func (s *TagService) Delete(id uuid.UUID) error {
	exists, err := s.tagRepository.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errx.ErrorNotFound(fmt.Errorf("tag with id %s is not found", id.String()))
	}

	existsByTea, err := s.tagRepository.ExistsByTeas(id)
	if existsByTea {
		return errx.ErrorBadRequest(fmt.Errorf("tag with id %s has some teas", id.String()))
	}

	err = s.tagRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
