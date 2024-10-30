package repository

import (
	"github.com/google/uuid"
	"repertoire/server/data/database"
	"repertoire/server/model"
)

type UserRepository interface {
	Get(user *model.User, id uuid.UUID) error
	GetByEmail(user *model.User, email string) error
	Create(user *model.User) error
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
	return u.client.DB.Find(&user, model.User{ID: id}).Error
}

func (u userRepository) GetByEmail(user *model.User, email string) error {
	return u.client.DB.Find(&user, model.User{Email: email}).Error
}

func (u userRepository) Create(user *model.User) error {
	return u.client.DB.Create(&user).Error
}
