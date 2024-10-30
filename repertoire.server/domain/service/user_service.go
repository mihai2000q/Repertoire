package service

import (
	"github.com/google/uuid"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"
)

type UserService interface {
	Get(id uuid.UUID) (user model.User, e *wrapper.ErrorCode)
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

func (s *userService) Get(id uuid.UUID) (model.User, *wrapper.ErrorCode) {
	return s.getUser.Handle(id)
}
