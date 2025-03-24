package service

import (
	"github.com/levchenki/tea-api/internal/entity"
)

type UserRepository interface {
	Exists(telegramId uint64) (bool, error)
	Create(user *entity.User) error
	GetByTelegramId(telegramId uint64) (*entity.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Create(user *entity.User) error {
	return s.userRepository.Create(user)
}

func (s *UserService) Exists(telegramId uint64) (bool, error) {
	return s.userRepository.Exists(telegramId)
}

func (s *UserService) GetByTelegramId(telegramId uint64) (*entity.User, error) {
	return s.userRepository.GetByTelegramId(telegramId)
}
