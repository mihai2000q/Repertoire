package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *models.Artist, id uuid.UUID) error
	GetAllByUser(artists *[]models.Artist, userId uuid.UUID) error
	Create(artist *models.Artist) error
	Update(artist *models.Artist) error
	Delete(id uuid.UUID) error
}

type artistRepository struct {
	client database.Client
}

func NewArtistRepository(client database.Client) ArtistRepository {
	return artistRepository{
		client: client,
	}
}

func (s artistRepository) Get(artist *models.Artist, id uuid.UUID) error {
	return s.client.DB.Find(&artist, models.Artist{ID: id}).Error
}

func (s artistRepository) GetAllByUser(artists *[]models.Artist, userId uuid.UUID) error {
	return s.client.DB.Model(&models.Artist{}).
		Where(models.Artist{UserID: userId}).
		Find(&artists).
		Error
}

func (s artistRepository) Create(artist *models.Artist) error {
	return s.client.DB.Create(&artist).Error
}

func (s artistRepository) Update(artist *models.Artist) error {
	return s.client.DB.Save(&artist).Error
}

func (s artistRepository) Delete(id uuid.UUID) error {
	return s.client.DB.Delete(&models.Artist{}, id).Error
}
