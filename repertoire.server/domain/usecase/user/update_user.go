package user

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateUser struct {
	repository repository.UserRepository
}

func NewUpdateUser(repository repository.UserRepository) UpdateUser {
	return UpdateUser{repository: repository}
}

func (u UpdateUser) Handle(request requests.UpdateUserRequest) *wrapper.ErrorCode {
	var user model.User
	err := u.repository.Get(&user, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return wrapper.NotFoundError(errors.New("user not found"))
	}

	user.Name = request.Name

	err = u.repository.Update(&user)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
