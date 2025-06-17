package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas/teaSchemas"
)

type TeaRepository interface {
	GetById(id uuid.UUID) (*entity.TeaWithRating, error)
	GetByIdWithUser(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error)
	GetAll(filters *teaSchemas.Filters) ([]entity.TeaWithRating, uint64, error)
	Create(inputTea *teaSchemas.RequestModel) (*entity.Tea, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, inputTea *teaSchemas.RequestModel, tagsToInsert, tagsToDelete []uuid.UUID) (*entity.Tea, error)
	Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) error
	Exists(id uuid.UUID) (bool, error)
	ExistsByName(existedId uuid.UUID, name string) (bool, error)

	GetMinServePrice() (float64, error)
	GetMaxServePrice() (float64, error)

	SetFavourite(id, userId uuid.UUID) error
	RemoveFavourite(id, userId uuid.UUID) error
}

type TeaTagRepository interface {
	GetByTeaId(teaId uuid.UUID) ([]entity.Tag, error)
	GetAllByTeaIds(teaIds []uuid.UUID) (map[uuid.UUID][]entity.Tag, error)
}

type TeaService struct {
	teaRepository TeaRepository
	tagRepository TeaTagRepository
}

func NewTeaService(
	teaRepository TeaRepository,
	tagRepository TeaTagRepository,
) *TeaService {
	return &TeaService{
		teaRepository: teaRepository,
		tagRepository: tagRepository,
	}
}

func (s *TeaService) GetTeaById(id uuid.UUID, userId uuid.UUID) (*entity.TeaWithRating, error) {
	var teaById *entity.TeaWithRating
	var err error
	if userId == uuid.Nil {
		teaById, err = s.teaRepository.GetById(id)
	} else {
		teaById, err = s.teaRepository.GetByIdWithUser(id, userId)
	}
	if err != nil {
		return nil, err
	}

	if teaById == nil {
		err := fmt.Errorf("tea with id %s is not found", id.String())
		return nil, errx.NewNotFoundError(err)
	}

	tags, err := s.tagRepository.GetByTeaId(id)
	if err != nil {
		return nil, err
	}
	teaById.Tags = tags

	return teaById, nil
}

func (s *TeaService) GetAllTeas(filters *teaSchemas.Filters) ([]entity.TeaWithRating, uint64, error) {
	allTeas, total, err := s.teaRepository.GetAll(filters)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 || len(allTeas) == 0 {
		return allTeas, 0, nil
	}

	teaIds := make([]uuid.UUID, len(allTeas))
	for i := range allTeas {
		teaIds[i] = allTeas[i].Id
	}

	tagsByTeaId, err := s.tagRepository.GetAllByTeaIds(teaIds)
	if err != nil {
		return nil, 0, err
	}

	for i, t := range allTeas {
		tags := tagsByTeaId[t.Id]
		allTeas[i].Tags = tags
	}

	return allTeas, total, err
}

func (s *TeaService) CreateTea(t *teaSchemas.RequestModel) (*entity.Tea, error) {
	exists, err := s.teaRepository.ExistsByName(uuid.Nil, t.Name)
	if err != nil {
		return nil, err
	}
	if exists == true {
		err := fmt.Errorf("tea with name %s has already existed", t.Name)
		return nil, errx.NewBadRequestError(err)
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

func (s *TeaService) DeleteTea(id uuid.UUID) error {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return err
	}
	if exists == false {
		err := fmt.Errorf("tea with id %s is not found", id.String())
		return errx.NewNotFoundError(err)
	}
	err = s.teaRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TeaService) UpdateTea(id uuid.UUID, t *teaSchemas.RequestModel) (*entity.Tea, error) {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists == false {
		err := fmt.Errorf("tea with id %s is not found", id.String())
		return nil, errx.NewNotFoundError(err)
	}

	existsByName, err := s.teaRepository.ExistsByName(id, t.Name)
	if err != nil {
		return nil, err
	}
	if existsByName == true {
		err := fmt.Errorf("tea with name %s has already existed", t.Name)
		return nil, errx.NewBadRequestError(err)
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

func (s *TeaService) getTagsDelta(existedTagIds, incomingTagIds []uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
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

func (s *TeaService) Evaluate(id uuid.UUID, userId uuid.UUID, evaluation *teaSchemas.Evaluation) (*entity.TeaWithRating, error) {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists == false {
		err := fmt.Errorf("tea with id %s is not found", id.String())
		return nil, errx.NewNotFoundError(err)
	}

	err = s.teaRepository.Evaluate(id, userId, evaluation)
	if err != nil {
		return nil, err
	}

	evaluatedTea, err := s.teaRepository.GetByIdWithUser(id, userId)
	if err != nil {
		return nil, err
	}
	return evaluatedTea, nil
}

func (s *TeaService) GetMinMaxServePrices() (float64, float64, error) {
	minPrice, err := s.teaRepository.GetMinServePrice()
	if err != nil {
		return 0, 0, err
	}
	maxPrice, err := s.teaRepository.GetMaxServePrice()
	if err != nil {
		return 0, 0, err
	}
	return minPrice, maxPrice, nil
}

func (s *TeaService) ToggleFavourites(id uuid.UUID, userId uuid.UUID, isFavourite bool) error {
	exists, err := s.teaRepository.Exists(id)
	if err != nil {
		return err
	}

	if !exists {
		err := fmt.Errorf("tea with id %s is not found", id.String())
		return errx.NewNotFoundError(err)
	}

	if isFavourite {
		err = s.teaRepository.SetFavourite(id, userId)
	} else {
		err = s.teaRepository.RemoveFavourite(id, userId)
	}
	if err != nil {
		return err
	}
	return nil
}
