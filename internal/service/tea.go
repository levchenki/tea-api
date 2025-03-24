package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
)

type TeaRepository interface {
	GetById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error)
	GetAll(filters *teaSchemas.Filters) ([]entity.TeaWithRating, error)
	Create(inputTea *teaSchemas.RequestModel) (*entity.Tea, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, inputTea *teaSchemas.RequestModel, tagsToInsert, tagsToDelete []uuid.UUID) (*entity.Tea, error)
	Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(existedId uuid.UUID, name string) (bool, error)
}

type TagRepository interface {
	GetByTeaId(teaId uuid.UUID) ([]entity.Tag, error)
}

type Service struct {
	teaRepository TeaRepository
	tagRepository TagRepository
}

func NewTeaService(
	teaRepository TeaRepository,
	tagRepository TagRepository,
) *Service {
	return &Service{
		teaRepository: teaRepository,
		tagRepository: tagRepository,
	}
}

func (s *Service) GetTeaById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error) {
	teaById, err := s.teaRepository.GetById(id, userId)
	if err != nil {
		return nil, err
	}

	if teaById == nil {
		return nil, errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
	}

	tags, err := s.tagRepository.GetByTeaId(id)
	if err != nil {
		return nil, err
	}
	teaById.Tags = tags

	return teaById, nil
}

func (s *Service) GetAllTeas(filters *teaSchemas.Filters, userId uuid.UUID) ([]entity.TeaWithRating, error) {
	filters.UserId = userId
	allTeas, err := s.teaRepository.GetAll(filters)
	if err != nil {
		return nil, err
	}

	return allTeas, err
}

func (s *Service) CreateTea(t *teaSchemas.RequestModel) (*entity.Tea, error) {
	exists, err := s.teaRepository.ExistsByName(uuid.Nil, t.Name)
	if err != nil {
		return nil, err
	}
	if exists == true {
		return nil, errx.ErrorBadRequest(fmt.Errorf("tea with name %s already is exist", t.Name))
	}

	createdTea, err := s.teaRepository.Create(t)
	if err != nil {
		return nil, err
	}

	tags, err := s.tagRepository.GetByTeaId(createdTea.Id)
	if err != nil {
		return nil, err
	}
	createdTea.Tags = tags

	return createdTea, nil
}

func (s *Service) DeleteTea(id uuid.UUID) error {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return err
	}
	if exists == false {
		return errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
	}
	err = s.teaRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateTea(id uuid.UUID, t *teaSchemas.RequestModel) (*entity.Tea, error) {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
	}

	existsByName, err := s.teaRepository.ExistsByName(id, t.Name)
	if err != nil {
		return nil, err
	}
	if existsByName == true {
		return nil, errx.ErrorBadRequest(fmt.Errorf("tea with name %s already is exist", t.Name))
	}

	tags, err := s.tagRepository.GetByTeaId(id)
	if err != nil {
		return nil, err
	}

	existedTagIds := make([]uuid.UUID, len(tags))
	for i, tag := range tags {
		existedTagIds[i] = tag.Id
	}

	tagsToInsert, tagsToDelete := s.getTagsDelta(existedTagIds, t.TagIds)

	updatedTea, err := s.teaRepository.Update(id, t, tagsToInsert, tagsToDelete)
	if err != nil {
		return nil, err
	}

	tags, err = s.tagRepository.GetByTeaId(id)
	updatedTea.Tags = tags
	return updatedTea, nil
}

func (s *Service) getTagsDelta(existedTagIds, incomingTagIds []uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
	existedTagsMap := make(map[uuid.UUID]uuid.UUID, len(existedTagIds))
	for _, tagId := range existedTagIds {
		existedTagsMap[tagId] = tagId
	}

	incomingTagsMap := make(map[uuid.UUID]uuid.UUID, len(incomingTagIds))
	for _, tagId := range incomingTagIds {
		incomingTagsMap[tagId] = tagId
	}

	tagIdsToDelete := make([]uuid.UUID, 0, len(existedTagIds))
	for _, tagId := range existedTagIds {
		if _, isOk := incomingTagsMap[tagId]; !isOk {
			tagIdsToDelete = append(tagIdsToDelete, tagId)
		}
	}

	tagIdsToInsert := make([]uuid.UUID, 0, len(incomingTagIds))
	for _, tagId := range incomingTagIds {
		if _, isOk := existedTagsMap[tagId]; !isOk {
			tagIdsToInsert = append(tagIdsToInsert, tagId)
		}
	}

	return tagIdsToInsert, tagIdsToDelete
}

func (s *Service) Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) (*entity.TeaWithRating, error) {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, errx.ErrorNotFound(fmt.Errorf("tea with id %s is not found", id.String()))
	}

	err = s.teaRepository.Evaluate(id, userId, evaluation)
	if err != nil {
		return nil, err
	}

	evaluatedTea, err := s.teaRepository.GetById(id, userId)
	if err != nil {
		return nil, err
	}
	return evaluatedTea, nil
}
