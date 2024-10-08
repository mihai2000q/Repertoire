package repository

import (
	"github.com/google/uuid"
	"repertoire/data/database"
	"repertoire/models"
)

type UserRepository struct {
	client database.Client
}

func NewUserRepository(client database.Client) UserRepository {
	return UserRepository{
		client: client,
	}
}

func (u UserRepository) Get(user *models.User, id uuid.UUID) error {
	return u.client.DB.Find(&user, models.User{ID: id}).Error
}

func (u UserRepository) GetByEmail(user *models.User, email string) error {
	return u.client.DB.Find(&user, models.User{Email: email}).Error
}

func (u UserRepository) Create(user *models.User) error {
	return u.client.DB.Create(&user).Error
}
