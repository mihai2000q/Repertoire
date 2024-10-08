package repository

import (
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

func (u UserRepository) GetByEmail(user *models.User, email string) error {
	// return u.client.DB.First(&user, "email = ?", email).Error
	return u.client.DB.Find(&user, models.User{Email: email}).Error
}

func (u UserRepository) Create(user *models.User) error {
	return u.client.DB.Create(&user).Error
}
