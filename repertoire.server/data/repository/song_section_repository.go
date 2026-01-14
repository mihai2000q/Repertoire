package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongSectionRepository interface {
	Get(section *model.SongSection, id uuid.UUID) error
	CountAllBySong(count *int64, songID uuid.UUID) error
	Create(section *model.SongSection) error
	Update(section *model.SongSection) error
	Delete(ids []uuid.UUID) error

	GetTypes(types *[]model.SongSectionType, userID uuid.UUID) error

	GetHistory(
		history *[]model.SongSectionHistory,
		sectionID uuid.UUID,
		property model.SongSectionProperty,
	) error
	CreateHistory(history *model.SongSectionHistory) error
}

type songSectionRepository struct {
	client database.Client
}

func NewSongSectionRepository(client database.Client) SongSectionRepository {
	return songSectionRepository{
		client: client,
	}
}

func (s songSectionRepository) Get(section *model.SongSection, id uuid.UUID) error {
	return s.client.Find(&section, model.SongSection{ID: id}).Error
}

func (s songSectionRepository) CountAllBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongSection{}).
		Where(model.SongSection{SongID: songID}).
		Count(count).
		Error
}

func (s songSectionRepository) Create(section *model.SongSection) error {
	return s.client.Create(&section).Error
}

func (s songSectionRepository) Update(section *model.SongSection) error {
	return s.client.Save(&section).Error
}

func (s songSectionRepository) Delete(ids []uuid.UUID) error {
	return s.client.Delete(&model.SongSection{}, ids).Error
}

// Types

func (s songSectionRepository) GetTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return s.client.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(&types).
		Error
}

// History

func (s songSectionRepository) GetHistory(
	history *[]model.SongSectionHistory,
	sectionID uuid.UUID,
	property model.SongSectionProperty,
) error {
	return s.client.
		Order("created_at").
		Find(&history, model.SongSectionHistory{SongSectionID: sectionID, Property: property}).
		Error
}

func (s songSectionRepository) CreateHistory(history *model.SongSectionHistory) error {
	return s.client.Create(&history).Error
}
