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

func (u SongRepository) Get(song *models.Song, id uuid.UUID) error {
	return u.client.DB.Find(&song, models.Song{ID: id}).Error
}

func (u SongRepository) Create(song *models.Song) error {
	return u.client.DB.Create(&song).Error
}
