package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	Get(album *models.Album, id uuid.UUID) error
	GetAllByUser(albums *[]models.Album, userId uuid.UUID) error
	Create(album *models.Album) error
	Update(album *models.Album) error
	Delete(id uuid.UUID) error
}

type albumRepository struct {
	client database.Client
}

func NewAlbumRepository(client database.Client) AlbumRepository {
	return albumRepository{
		client: client,
	}
}

func (a albumRepository) Get(album *models.Album, id uuid.UUID) error {
	return a.client.DB.Find(&album, models.Album{ID: id}).Error
}

func (a albumRepository) GetAllByUser(albums *[]models.Album, userId uuid.UUID) error {
	return a.client.DB.Model(&models.Album{}).
		Where(models.Album{UserID: userId}).
		Find(&albums).
		Error
}

func (a albumRepository) Create(album *models.Album) error {
	return a.client.DB.Create(&album).Error
}

func (a albumRepository) Update(album *models.Album) error {
	return a.client.DB.Save(&album).Error
}

func (a albumRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&models.Album{}, id).Error
}
