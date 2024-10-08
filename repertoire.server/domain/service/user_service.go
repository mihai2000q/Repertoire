package service

import (
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/models"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

func (s *UserService) Get(id uuid.UUID) (user models.User, err error) {
	return user, s.repository.Get(&user, id)
}
