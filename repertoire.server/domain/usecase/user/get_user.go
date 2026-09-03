package user

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type GetUser struct {
	repository repository.UserRepository
}

func NewGetUser(repository repository.UserRepository) GetUser {
	return GetUser{repository: repository}
}

func (g GetUser) Handle(id uuid.UUID) (user model.User, e *httperror.ErrorCode) {
	if err := g.repository.Get(&user, id); err != nil {
		return user, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return user, httperror.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
