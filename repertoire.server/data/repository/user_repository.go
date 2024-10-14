package repository

import (
	"github.com/google/uuid"
	"repertoire/data/database"
	"repertoire/models"
)

type UserRepository interface {
	Get(user *models.User, id uuid.UUID) error
	GetByEmail(user *models.User, email string) error
	Create(user *models.User) error
}

type userRepository struct {
	client database.Client
}

func NewUserRepository(client database.Client) UserRepository {
	return userRepository{
		client: client,
	}
}

func (u userRepository) Get(user *models.User, id uuid.UUID) error {
	return u.client.DB.Find(&user, models.User{ID: id}).Error
}

func (u userRepository) GetByEmail(user *models.User, email string) error {
	return u.client.DB.Find(&user, models.User{Email: email}).Error
}

func (u userRepository) Create(user *models.User) error {
	return u.client.DB.Create(&user).Error
}
