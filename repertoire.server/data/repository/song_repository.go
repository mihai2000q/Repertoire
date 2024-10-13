package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type SongRepository struct {
	client database.Client
}

func NewSongRepository(client database.Client) SongRepository {
	return SongRepository{
		client: client,
	}
}

func (s SongRepository) Get(song *models.Song, id uuid.UUID) error {
	return s.client.DB.Find(&song, models.Song{ID: id}).Error
}

func (s SongRepository) GetAllByUser(songs []models.Song, userId uuid.UUID) error {
	return s.client.DB.Model(&models.Song{}).Where(models.Song{ID: userId}).Find(&songs).Error
}

func (s SongRepository) Create(song *models.Song) error {
	return s.client.DB.Create(&song).Error
}

func (s SongRepository) Update(song *models.Song) error {
	return s.client.DB.Save(&song).Error
}

func (s SongRepository) Delete(song *models.Song) error {
	return s.client.DB.Delete(&song).Error
}
