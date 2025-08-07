package user

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type GetUser struct {
	repository repository.UserRepository
}

func NewGetUser(repository repository.UserRepository) GetUser {
	return GetUser{repository: repository}
}

func (g GetUser) Handle(id uuid.UUID) (user model.User, e *wrapper.ErrorCode) {
	err := g.repository.Get(&user, id)
	if err != nil {
		return user, wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return user, wrapper.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
