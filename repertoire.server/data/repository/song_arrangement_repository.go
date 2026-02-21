package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongArrangementRepository interface {
	GetWithAssociations(arrangement *model.SongArrangement, id uuid.UUID) error
	GetAllBySong(arrangements *[]model.SongArrangement, songID uuid.UUID) error
	CountBySong(count *int64, songID uuid.UUID) error
	Create(arrangement *model.SongArrangement) error
	UpdateWithAssociations(arrangement *model.SongArrangement) error
	UpdateAllWithAssociations(arrangements *[]model.SongArrangement) error
	Delete(id uuid.UUID) error
}

type songArrangementRepository struct {
	client database.Client
}

func NewSongArrangementRepository(client database.Client) SongArrangementRepository {
	return songArrangementRepository{
		client: client,
	}
}

func (s songArrangementRepository) GetWithAssociations(arrangement *model.SongArrangement, id uuid.UUID) error {
	return s.client.
		Preload("SectionOccurrences").
		Find(&arrangement, model.SongArrangement{ID: id}).
		Error
}

func (s songArrangementRepository) GetAllBySong(arrangements *[]model.SongArrangement, songID uuid.UUID) error {
	return s.client.Model(&model.SongArrangement{}).
		Preload("SectionOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_sections ON song_sections.song_id = song_arrangements.song_id").
				Order("song_sections.order")
		}).
		Preload("SectionOccurrences.Section"). // Try Joins
		Where(model.SongArrangement{SongID: songID}).
		Order("\"order\"").
		Find(arrangements).
		Error
}

func (s songArrangementRepository) CountBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongArrangement{}).
		Where("song_id = ?", songID).
		Count(count).
		Error
}

func (s songArrangementRepository) Create(arrangement *model.SongArrangement) error {
	return s.client.Create(&arrangement).Error
}

func (s songArrangementRepository) UpdateWithAssociations(arrangement *model.SongArrangement) error {
	return s.client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&arrangement).
		Error
}

func (s songArrangementRepository) UpdateAllWithAssociations(arrangements *[]model.SongArrangement) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, arrangement := range *arrangements {
			err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&arrangement).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songArrangementRepository) Delete(id uuid.UUID) error {
	return s.client.Delete(&model.SongArrangement{}, id).Error
}
