package repository

import (
	"repertoire/data/database"
	"repertoire/model"

	"github.com/google/uuid"
)

type SongRepository interface {
	Get(song *model.Song, id uuid.UUID) error
	GetAllByUser(
		songs *[]model.Song,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID) error
	GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error
	Create(song *model.Song) error
	Update(song *model.Song) error
	Delete(id uuid.UUID) error
}

type songRepository struct {
	client database.Client
}

func NewSongRepository(client database.Client) SongRepository {
	return songRepository{
		client: client,
	}
}

func (s songRepository) Get(song *model.Song, id uuid.UUID) error {
	return s.client.DB.Find(&song, model.Song{ID: id}).Error
}

func (s songRepository) GetAllByUser(
	songs *[]model.Song,
	userID uuid.UUID,
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
	return s.client.DB.Model(&model.Song{}).
		Where(model.Song{UserID: userID}).
		Order(orderBy).
		Offset(offset).
		Limit(*pageSize).
		Find(&songs).
		Error
}

func (s songRepository) GetAllByUserCount(count *int64, userID uuid.UUID) error {
	return s.client.DB.Model(&model.Song{}).
		Where(model.Song{UserID: userID}).
		Count(count).
		Error
}

func (s songRepository) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	return s.client.DB.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Find(&tunings).
		Error
}

func (s songRepository) Create(song *model.Song) error {
	return s.client.DB.Create(&song).Error
}

func (s songRepository) Update(song *model.Song) error {
	return s.client.DB.Save(&song).Error
}

func (s songRepository) Delete(id uuid.UUID) error {
	return s.client.DB.Delete(&model.Song{}, id).Error
}
