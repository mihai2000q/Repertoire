package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UserService interface {
	Get(id uuid.UUID) (user model.User, e *wrapper.ErrorCode)
	Update(request requests.UpdateUserRequest) *wrapper.ErrorCode
}

type userService struct {
	getUser    user.GetUser
	updateUser user.UpdateUser
}

func NewUserService(
	getUser user.GetUser,
	updateUser user.UpdateUser,
) UserService {
	return &userService{
		getUser:    getUser,
		updateUser: updateUser,
	}
}

func (s *userService) Get(id uuid.UUID) (model.User, *wrapper.ErrorCode) {
	return s.getUser.Handle(id)
}

func (s *userService) Update(request requests.UpdateUserRequest) *wrapper.ErrorCode {
	return s.updateUser.Handle(request)
}
