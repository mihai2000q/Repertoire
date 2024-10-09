package service

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

func (s *UserService) Get(id uuid.UUID) (user models.User, e *utils.ErrorCode) {
	err := s.repository.Get(&user, id)
	if err != nil {
		return user, utils.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return user, utils.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
