package repository

import (
	"repertoire/data/database"
	"repertoire/models"
)

type SongRepository struct {
	client database.Client
}

func NewSongRepository(client database.Client) SongRepository {
	return SongRepository{
		client: client,
	}
}

func (u SongRepository) Create(song *models.Song) error {
	return u.client.DB.Create(&song).Error
}
