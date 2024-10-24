package repository

import (
	"repertoire/data/database"
	"repertoire/model"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *model.Artist, id uuid.UUID) error
	GetAllByUser(
		artists *[]model.Artist,
		userId uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy string,
	) error
	GetAllByUserCount(count *int64, userId uuid.UUID) error
	Create(artist *model.Artist) error
	Update(artist *model.Artist) error
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

func (a artistRepository) Get(artist *model.Artist, id uuid.UUID) error {
	return a.client.DB.Find(&artist, model.Artist{ID: id}).Error
}

func (a artistRepository) GetAllByUser(
	artists *[]model.Artist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy string,
) error {
	offset := -1
	if pageSize == nil {
		pageSize = &[]int{-1}[0]
	} else {
		offset = (*currentPage - 1) * *pageSize
	}
	return a.client.DB.Model(&model.Artist{}).
		Where(model.Artist{UserID: userId}).
		Order(orderBy).
		Offset(offset).
		Limit(*pageSize).
		Find(&artists).
		Error
}

func (a artistRepository) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	return a.client.DB.Model(&model.Artist{}).
		Where(model.Artist{UserID: userId}).
		Count(count).
		Error
}

func (a artistRepository) Create(artist *model.Artist) error {
	return a.client.DB.Create(&artist).Error
}

func (a artistRepository) Update(artist *model.Artist) error {
	return a.client.DB.Save(&artist).Error
}

func (a artistRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&model.Artist{}, id).Error
}
