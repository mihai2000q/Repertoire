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
	userRepository repository.UserRepository
}

func NewGetUser(userRepository repository.UserRepository) GetUser {
	return GetUser{userRepository: userRepository}
}

func (g GetUser) Handle(id uuid.UUID) (user model.User, e *httperror.ErrorCode) {
	if err := g.userRepository.Get(&user, id); err != nil {
		return user, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return user, httperror.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
