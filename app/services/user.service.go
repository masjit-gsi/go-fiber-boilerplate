package services

import (
	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/app/repository"
)

type UserService interface {
	GetUserByID(id string) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: repository,
	}
}

func (s *UserServiceImpl) GetUserByID(id string) (user models.User, err error) {
	return s.UserRepository.GetUserByID(id)
}

func (s *UserServiceImpl) GetUserByUsername(username string) (user models.User, err error) {
	return s.UserRepository.GetUserByUsername(username)
}
