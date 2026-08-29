package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongPartRepository interface {
	Get(instrumentPart *model.SongPart, id uuid.UUID) error
	CountAllBySection(count *int64, sectionID uuid.UUID) error
	CountAllBySong(count *int64, songID uuid.UUID) error
	Create(instrumentPart *model.SongPart) error
	Update(instrumentPart *model.SongPart) error
	Delete(ids []uuid.UUID) error

	GetHistory(
		history *[]model.SongPartHistory,
		instrumentPartID uuid.UUID,
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

func (s songPartRepository) Get(instrumentPart *model.SongPart, id uuid.UUID) error {
	return s.client.Find(&instrumentPart, model.SongPart{ID: id}).Error
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

func (s songPartRepository) Create(instrumentPart *model.SongPart) error {
	return s.client.Create(&instrumentPart).Error
}

func (s songPartRepository) Update(instrumentPart *model.SongPart) error {
	return s.client.Save(&instrumentPart).Error
}

func (s songPartRepository) Delete(ids []uuid.UUID) error {
	return s.client.Delete(&model.SongPart{}, ids).Error
}

// History

func (s songPartRepository) GetHistory(
	history *[]model.SongPartHistory,
	instrumentPartID uuid.UUID,
	property model.SongPartProperty,
) error {
	return s.client.
		Order("created_at").
		Find(&history, model.SongPartHistory{SongPartID: instrumentPartID, Property: property}).
		Error
}

func (s songPartRepository) CreateHistory(history *model.SongPartHistory) error {
	return s.client.Create(&history).Error
}
