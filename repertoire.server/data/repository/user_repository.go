package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Get(user *model.User, id uuid.UUID) error
	GetByEmail(user *model.User, email string) error
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uuid.UUID) error
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

func (u userRepository) Update(user *model.User) error {
	return u.client.DB.Save(&user).Error
}

func (u userRepository) Delete(id uuid.UUID) error {
	return u.client.DB.Delete(&model.User{}, id).Error
}
