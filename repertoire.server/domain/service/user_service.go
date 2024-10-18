package service

import (
	"github.com/google/uuid"
	"repertoire/domain/usecases/user"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type UserService interface {
	Get(id uuid.UUID) (user models.User, e *wrapper.ErrorCode)
}

type userService struct {
	getUser user.GetUser
}

func NewUserService(
	getUser user.GetUser,
) UserService {
	return &userService{
		getUser: getUser,
	}
}

func (s *userService) Get(id uuid.UUID) (models.User, *wrapper.ErrorCode) {
	return s.getUser.Handle(id)
}
