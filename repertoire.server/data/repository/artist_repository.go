package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *models.Artist, id uuid.UUID) error
	GetAllByUser(artists *[]models.Artist, userId uuid.UUID, currentPage *int, pageSize *int) error
	GetAllByUserCount(count *int64, userId uuid.UUID) error
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

func (a artistRepository) Get(artist *models.Artist, id uuid.UUID) error {
	return a.client.DB.Find(&artist, models.Artist{ID: id}).Error
}

func (a artistRepository) GetAllByUser(
	artists *[]models.Artist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
) error {
	return a.client.DB.Model(&models.Artist{}).
		Where(models.Artist{UserID: userId}).
		Offset((*currentPage - 1) * *pageSize).
		Limit(*pageSize).
		Find(&artists).
		Error
}

func (a artistRepository) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	return a.client.DB.Model(&models.Artist{}).
		Where(models.Artist{UserID: userId}).
		Count(count).
		Error
}

func (a artistRepository) Create(artist *models.Artist) error {
	return a.client.DB.Create(&artist).Error
}

func (a artistRepository) Update(artist *models.Artist) error {
	return a.client.DB.Save(&artist).Error
}

func (a artistRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&models.Artist{}, id).Error
}
