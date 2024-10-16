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

func (s albumRepository) Get(album *models.Album, id uuid.UUID) error {
	return s.client.DB.Find(&album, models.Album{ID: id}).Error
}

func (s albumRepository) GetAllByUser(albums *[]models.Album, userId uuid.UUID) error {
	return s.client.DB.Model(&models.Album{}).
		Where(models.Album{UserID: userId}).
		Find(&albums).
		Error
}

func (s albumRepository) Create(album *models.Album) error {
	return s.client.DB.Create(&album).Error
}

func (s albumRepository) Update(album *models.Album) error {
	return s.client.DB.Save(&album).Error
}

func (s albumRepository) Delete(id uuid.UUID) error {
	return s.client.DB.Delete(&models.Album{}, id).Error
}
