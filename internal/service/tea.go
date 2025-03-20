package service

import (
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
)

type TeaRepository interface {
	GetById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error)
	GetAll(filters *teaSchemas.Filters) ([]entity.TeaWithRating, error)
	Create(inputTea *teaSchemas.RequestModel) (*entity.Tea, error)
	Delete(id uuid.UUID) (bool, error)
	Update(id uuid.UUID, inputTea *teaSchemas.RequestModel, tagsToInsert, tagsToDelete []uuid.UUID) (*entity.Tea, error)
}

type TagRepository interface {
	GetByTeaId(teaId uuid.UUID) ([]entity.Tag, error)
}

type UserRepository interface {
	GetId(telegramId int) (uuid.UUID, error)
}

type Service struct {
	teaRepository  TeaRepository
	tagRepository  TagRepository
	userRepository UserRepository
}

func NewTeaService(
	teaRepository TeaRepository,
	tagRepository TagRepository,
	userRepository UserRepository,
) *Service {
	return &Service{
		teaRepository:  teaRepository,
		tagRepository:  tagRepository,
		userRepository: userRepository,
	}
}

func (s *Service) GetTeaById(id uuid.UUID, telegramUserId int) (*entity.TeaWithRating, error) {
	userId, err := s.userRepository.GetId(telegramUserId)
	if err != nil {
		return nil, err
	}

	teaById, err := s.teaRepository.GetById(id, userId)
	if err != nil {
		return nil, err
	}

	if teaById == nil {
		return nil, nil
	}

	tags, err := s.tagRepository.GetByTeaId(id)
	if err != nil {
		return nil, err
	}
	teaById.Tags = tags

	return teaById, nil
}

func (s *Service) GetAllTeas(filters *teaSchemas.Filters, telegramUserId int) ([]entity.TeaWithRating, error) {
	userId, err := s.userRepository.GetId(telegramUserId)
	if err != nil {
		return nil, err
	}
	filters.UserId = userId
	allTeas, err := s.teaRepository.GetAll(filters)
	if err != nil {
		return nil, err
	}

	return allTeas, err
}

func (s *Service) CreateTea(t *teaSchemas.RequestModel) (*entity.Tea, error) {
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

func (s *Service) DeleteTea(id uuid.UUID) (bool, error) {
	isDeleted, err := s.teaRepository.Delete(id)
	if err != nil {
		return false, err
	}
	return isDeleted, nil
}

func (s *Service) UpdateTea(id uuid.UUID, t *teaSchemas.RequestModel) (*entity.Tea, error) {
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

	if updatedTea == nil {
		return nil, nil
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
