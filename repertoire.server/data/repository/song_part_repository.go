package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongPartRepository interface {
	Get(part *model.SongPart, id uuid.UUID) error
	GetAllByIDs(parts *[]model.SongPart, ids []uuid.UUID) error
	CountAllBySong(count *int64, songID uuid.UUID) error
	CountBySectionIDs(sectionIDs []uuid.UUID) (map[uuid.UUID]int64, error)
	Create(part *model.SongPart) error
	Update(part *model.SongPart) error
	UpdateAll(parts *[]model.SongPart) error
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

func (s songPartRepository) GetAllByIDs(parts *[]model.SongPart, ids []uuid.UUID) error {
	return s.client.Find(&parts, ids).Error
}

func (s songPartRepository) CountAllBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongPart{}).
		Where(model.SongPart{SongID: songID}).
		Count(count).
		Error
}

func (s songPartRepository) CountBySectionIDs(sectionIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	type result struct {
		SectionID uuid.UUID
		Count     int64
	}
	var results []result
	err := s.client.Model(&model.SongSectionPart{}).
		Select("section_id, COUNT(*) AS count").
		Where("section_id IN ?", sectionIDs).
		Group("section_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int64)
	for _, res := range results {
		counts[res.SectionID] = res.Count
	}
	return counts, nil
}

func (s songPartRepository) Create(part *model.SongPart) error {
	return s.client.Create(&part).Error
}

func (s songPartRepository) Update(part *model.SongPart) error {
	return s.client.Save(&part).Error
}

func (s songPartRepository) UpdateAll(parts *[]model.SongPart) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, part := range *parts {
			if err := tx.Save(&part).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
		Find(&history, model.SongPartHistory{PartID: partID, Property: property}).
		Error
}

func (s songPartRepository) CreateHistory(history *model.SongPartHistory) error {
	return s.client.Create(&history).Error
}
