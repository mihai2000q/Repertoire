package user

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetUser struct {
	repository repository.UserRepository
}

func NewGetUser(repository repository.UserRepository) GetUser {
	return GetUser{repository: repository}
}

func (g GetUser) Handle(id uuid.UUID) (user models.User, e *utils.ErrorCode) {
	err := g.repository.Get(&user, id)
	if err != nil {
		return user, utils.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return user, utils.NotFoundError(errors.New("user not found"))
	}
	return user, nil
}
