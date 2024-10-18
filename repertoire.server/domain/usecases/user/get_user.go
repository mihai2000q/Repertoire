package user

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetUser struct {
	repository repository.UserRepository
}

func NewGetUser(repository repository.UserRepository) GetUser {
	return GetUser{repository: repository}
}

func (g GetUser) Handle(id uuid.UUID) (user models.User, e *wrapper.ErrorCode) {
	err := g.repository.Get(&user, id)
	if err != nil {
		return user, wrapper.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return user, wrapper.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
