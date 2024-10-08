package service

import (
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

func (s *UserService) GetByEmail(email string) (user models.User, err error) {
	return user, s.repository.GetByEmail(&user, email)
}
