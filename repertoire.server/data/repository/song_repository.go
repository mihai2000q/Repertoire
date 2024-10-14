package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type SongRepository interface {
	Get(song *models.Song, id uuid.UUID) error
	GetAllByUser(songs *[]models.Song, userId uuid.UUID) error
	Create(song *models.Song) error
	Update(song *models.Song) error
	Delete(song *models.Song) error
}

type songRepository struct {
	client database.Client
}

func NewSongRepository(client database.Client) SongRepository {
	return songRepository{
		client: client,
	}
}

func (s songRepository) Get(song *models.Song, id uuid.UUID) error {
	return s.client.DB.Find(&song, models.Song{ID: id}).Error
}

func (s songRepository) GetAllByUser(songs *[]models.Song, userId uuid.UUID) error {
	return s.client.DB.Model(&models.Song{}).
		Where(models.Song{UserID: userId}).
		Find(&songs).
		Error
}

func (s songRepository) Create(song *models.Song) error {
	return s.client.DB.Create(&song).Error
}

func (s songRepository) Update(song *models.Song) error {
	return s.client.DB.Save(&song).Error
}

func (s songRepository) Delete(song *models.Song) error {
	return s.client.DB.Delete(&song).Error
}
