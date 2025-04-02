package repository

import (
	"github.com/google/uuid"
	"repertoire/auth/data/database"
	"repertoire/auth/model"
)

type UserRepository interface {
	Get(user *model.User, id uuid.UUID) error
	GetByEmail(user *model.User, email string) error
}

type userRepository struct {
	client database.Client
}

func NewUserRepository(client database.Client) UserRepository {
	return userRepository{
		client: client,
	}
}

func (u userRepository) Get(user *model.User, id uuid.UUID) error {
	return u.client.Find(&user, model.User{ID: id}).Error
}

func (u userRepository) GetByEmail(user *model.User, email string) error {
	return u.client.Find(&user, model.User{Email: email}).Error
}
