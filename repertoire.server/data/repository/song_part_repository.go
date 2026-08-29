package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongPartRepository interface {
	Get(part *model.SongPart, id uuid.UUID) error
	CountAllBySection(count *int64, sectionID uuid.UUID) error
	CountAllBySong(count *int64, songID uuid.UUID) error
	Create(part *model.SongPart) error
	Update(part *model.SongPart) error
	Delete(ids []uuid.UUID) error

	GetHistory(
		history *[]model.SongPartHistory,
		partID uuid.UUID,
		property model.SongPartProperty,
	) error
	CreateHistory(history *model.SongPartHistory) error
}

type songPartRepository struct {
	client database.Client
}

func NewSongPartRepository(client database.Client) SongPartRepository {
	return songPartRepository{
		client: client,
	}
}

func (s songPartRepository) Get(part *model.SongPart, id uuid.UUID) error {
	return s.client.Find(&part, model.SongPart{ID: id}).Error
}

func (s songPartRepository) CountAllBySection(count *int64, sectionID uuid.UUID) error {
	return s.client.Model(&model.SongPart{}).
		Where(model.SongPart{SectionID: &sectionID}).
		Count(count).
		Error
}

func (s songPartRepository) CountAllBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongPart{}).
		Where(model.SongPart{SongID: songID}).
		Count(count).
		Error
}

func (s songPartRepository) Create(part *model.SongPart) error {
	return s.client.Create(&part).Error
}

func (s songPartRepository) Update(part *model.SongPart) error {
	return s.client.Save(&part).Error
}

func (s songPartRepository) Delete(ids []uuid.UUID) error {
	return s.client.Delete(&model.SongPart{}, ids).Error
}

// History

func (s songPartRepository) GetHistory(
	history *[]model.SongPartHistory,
	partID uuid.UUID,
	property model.SongPartProperty,
) error {
	return s.client.
		Order("created_at").
		Find(&history, model.SongPartHistory{SongPartID: partID, Property: property}).
		Error
}

func (s songPartRepository) CreateHistory(history *model.SongPartHistory) error {
	return s.client.Create(&history).Error
}
