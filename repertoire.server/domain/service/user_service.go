package service

import (
	"github.com/google/uuid"
	"repertoire/domain/usecases/user"
	"repertoire/models"
	"repertoire/utils"
)

type UserService interface {
	Get(id uuid.UUID) (user models.User, e *utils.ErrorCode)
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

func (s *userService) Get(id uuid.UUID) (models.User, *utils.ErrorCode) {
	return s.getUser.Handle(id)
}
